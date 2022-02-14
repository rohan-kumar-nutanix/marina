//
// Copyright (c) 2017 Nutanix Inc. All rights reserved.

// Author: zxie@nutanix.com

// Package to obtain a JWT token that can be used to access v3 REST API going through
// the authentication and access control layers of aplos.

package aplos_transport

import (
	"crypto/x509"
	"encoding/pem"
	"flag"
	"os/exec"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/errors"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	net "github.com/nutanix-core/acs-aos-go/nutanix/util-slbufs/util/sl_bufs/net"
	nutanix_prism "github.com/nutanix-core/acs-aos-go/prism_secure_key_repo"
	"github.com/nutanix-core/acs-aos-go/zeus"
)

const (
	ZK_KEY_REPO = "/appliance/logical/prismsecurekey"
)

var (
	privateKeyReadCmd = flag.String("go_jwt_private_key_read_sh_command",
		"/bin/cat",
		"Shell command used to read private key.")
)

// Obtain a JWT token to be used for v3 REST API call. This will fallback to
// aplos key in case of ecdsa. Ref- ENG-233390
func GetJWTWithFallback(zkSession *zeus.ZookeeperSession,
	userUuid *uuid4.Uuid, tenantUuid *uuid4.Uuid,
	userName *string, reqContext *net.RpcRequestContext) (string, error) {

	if zkSession == nil {
		glog.Error("ZK session is nil")
		return "", nil
	}

	zkBytes, _, err := zkSession.Conn.GetWithBackoff(ZK_KEY_REPO)
	if err != nil {
		return "", err
	}
	keyRepo := &nutanix_prism.SecureKeyRepository{}
	err = proto.Unmarshal(zkBytes, keyRepo)
	if err != nil {
		return "", err
	}

	// Similar to aplos logic, in case of ecdsa algorithm aplos key will be used
	// Added fallback logic like server_cert_utils. Ref- CALM-12549
	var keyEntry *nutanix_prism.SecureKeyRepository_SecureKey = keyRepo.GetClusterGeneratedKey()
	if len(keyRepo.GetPrismSecureKeyList()) > 0 &&
		keyRepo.GetPrismSecureKeyList()[0].GetEnabled() {
		keyEntry = keyRepo.GetPrismSecureKeyList()[0]
	}
	if keyEntry.GetType() == nutanix_prism.SecureKeyRepository_SecureKey_RSA_2048 {
		return GetJWT(zkSession, userUuid, tenantUuid, false, userName, reqContext)
	} else if keyEntry.GetType() == nutanix_prism.SecureKeyRepository_SecureKey_ECDSA_256 ||
		keyEntry.GetType() == nutanix_prism.SecureKeyRepository_SecureKey_ECDSA_384 {
		return GetJWT(zkSession, userUuid, tenantUuid, true, userName, reqContext)
	} else {
		return "", errors.NewNtnxError("Unsupported key type: "+
			keyEntry.GetType().String(), errors.ConfigErrorType)
	}
}

// Obtain a JWT token to be used for v3 REST API call. Using the provided
// user UUID and tenant UUID, if the UUID is nil, the default magic vaule
// of all zero UUID is used. Username can be passed if available.
func GetJWT(zkSession *zeus.ZookeeperSession,
	userUuid *uuid4.Uuid, tenantUuid *uuid4.Uuid,
	useAplosKey bool, userName *string, reqContext *net.RpcRequestContext) (string, error) {
	if zkSession == nil {
		glog.Error("ZK session is nil")
		return "", nil
	}

	zkBytes, _, err := zkSession.Conn.GetWithBackoff(ZK_KEY_REPO)
	if err != nil {
		return "", err
	}
	keyRepo := &nutanix_prism.SecureKeyRepository{}
	err = proto.Unmarshal(zkBytes, keyRepo)
	if err != nil {
		return "", err
	}

	// Mirror the aplos logic of choosing the default generated private key
	// if no enabled user provided private key found.
	var keyEntry *nutanix_prism.SecureKeyRepository_SecureKey = keyRepo.GetClusterGeneratedKey()
	if useAplosKey {
		keyEntry = keyRepo.GetAplosSecureKey()
	} else if len(keyRepo.GetPrismSecureKeyList()) > 0 &&
		keyRepo.GetPrismSecureKeyList()[0].GetEnabled() {
		keyEntry = keyRepo.GetPrismSecureKeyList()[0]
	}

	keyContent, err := exec.Command("/bin/sh", "-c", *privateKeyReadCmd+" "+
		keyEntry.GetPrivateKeyLocation()).Output()
	if err != nil {
		return "", err
	}
	keyContentDecoded, _ := pem.Decode(keyContent)
	if keyContentDecoded == nil {
		glog.Error("no pem block found")
		return "", errors.NewNtnxError("Failed to decode pem key",
			errors.ConfigErrorType)
	}

	var signingMethod jwt.SigningMethod
	var key interface{}
	if keyEntry.GetType() == nutanix_prism.SecureKeyRepository_SecureKey_RSA_2048 {
		signingMethod = jwt.SigningMethodRS256
		key, err = x509.ParsePKCS8PrivateKey(keyContentDecoded.Bytes)
		// ENG-124125 we need to parse ANS1 based key differently (PKCS1)
		if err != nil {
			glog.Info("Parsing aplos specific key: ", err)
			key, err = x509.ParsePKCS1PrivateKey(keyContentDecoded.Bytes)
		}
	} else if keyEntry.GetType() == nutanix_prism.SecureKeyRepository_SecureKey_ECDSA_256 {
		signingMethod = jwt.SigningMethodES256
		key, err = x509.ParseECPrivateKey(keyContentDecoded.Bytes)
	} else if keyEntry.GetType() == nutanix_prism.SecureKeyRepository_SecureKey_ECDSA_384 {
		signingMethod = jwt.SigningMethodES384
		key, err = x509.ParseECPrivateKey(keyContentDecoded.Bytes)
	} else {
		return "", errors.NewNtnxError("Unsupported key type: "+
			keyEntry.GetType().String(), errors.ConfigErrorType)
	}
	if err != nil {
		glog.Error("Failed to parse key content", err)
		return "", err
	}

	userUuidStr := "00000000-0000-0000-0000-000000000000"
	if userUuid != nil {
		userUuidStr = userUuid.String()
	}

	tenantUuidStr := "00000000-0000-0000-0000-000000000000"
	if tenantUuid != nil {
		tenantUuidStr = tenantUuid.String()
	}

	claims := jwt.MapClaims{
		"iss":         "aplos:local_service",
		"aud":         "aplos:cluster_services",
		"user_uuid":   userUuidStr,
		"tenant_uuid": tenantUuidStr,
	}

	if userName != nil {
		claims["username"] = *userName
	}

	if reqContext != nil {
		var serviceName = reqContext.GetServiceName()
		if serviceName != "" {
			claims["service_name"] = serviceName
		}
	}

	token := jwt.NewWithClaims(signingMethod, claims)

	return token.SignedString(key)
}
