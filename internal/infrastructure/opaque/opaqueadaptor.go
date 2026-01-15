package opaque

import (
	"os"

	"github.com/TheAmirhosssein/cool-password-manage/config"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/log"
	"github.com/bytemare/opaque"
)

type opaqueAdaptor struct {
	server *opaque.Server
	config *config.Config
}

func New(config *config.Config) (OpaqueService, error) {
	a := &opaqueAdaptor{config: config}
	return a, a.Init()
}

func (o *opaqueAdaptor) Init() error {
	serverID := []byte(o.config.Opaque.ServerID)

	conf := opaque.DefaultConfiguration()

	secretOprfSeed, err := o.getOprfKey()
	if err != nil {
		log.ErrorLogger.Error("error at getting oprf key", "error", err.Error())
		return err
	}

	serverPrivateKey, err := o.getPrivateKey()
	if err != nil {
		log.ErrorLogger.Error("error at getting private key", "error", err.Error())
		return err
	}

	serverPublicKey, err := o.getPublicKey()
	if err != nil {
		log.ErrorLogger.Error("error at getting public key", "error", err.Error())
		return err
	}

	server, err := conf.Server()
	if err != nil {
		log.ErrorLogger.Error("error at starting opaque server", "error", err.Error())
		return err
	}

	if err := server.SetKeyMaterial(serverID, serverPrivateKey, serverPublicKey, secretOprfSeed); err != nil {
		log.ErrorLogger.Error("error at setting key material", "error", err.Error())
		return err
	}

	o.server = server
	return nil
}

func (o *opaqueAdaptor) RegisterInit(message []byte) ([]byte, []byte, error) {
	req, err := o.server.Deserialize.RegistrationRequest(message)
	if err != nil {
		log.ErrorLogger.Error("error at deserializing message", "error", err.Error())
		return nil, nil, err
	}

	credID := opaque.RandomBytes(64)

	serverPublicKey, err := o.getPublicKey()
	if err != nil {
		log.ErrorLogger.Error("error at getting public key", "error", err.Error())
		return nil, nil, err
	}

	pks, err := o.server.Deserialize.DecodeAkePublicKey(serverPublicKey)
	if err != nil {
		log.ErrorLogger.Error("error at decoding ake public key", "error", err.Error())
		return nil, nil, err
	}

	secretOprfKey, err := o.getOprfKey()
	if err != nil {
		log.ErrorLogger.Error("error at getting oprf key", "error", err.Error())
		return nil, nil, err
	}

	resp := o.server.RegistrationResponse(req, pks, credID, secretOprfKey)

	return resp.Serialize(), credID, nil
}

func (o *opaqueAdaptor) RegisterFinalize(message, credID []byte, username string) ([]byte, error) {
	record, err := o.server.Deserialize.RegistrationRecord(message)
	if err != nil {
		log.ErrorLogger.Error("error at deserializing message", "error", err.Error())
		return nil, err
	}

	clientRecord := &opaque.ClientRecord{
		CredentialIdentifier: credID,           // used during serialization
		ClientIdentity:       []byte(username), // used during serialization
		RegistrationRecord:   record,
	}

	return clientRecord.Serialize(), nil
}

func (o *opaqueAdaptor) LoginInit(message []byte, opaqueRecord *opaque.ClientRecord) ([]byte, error) {
	ke1, err := o.server.Deserialize.KE1(message)
	if err != nil {
		log.ErrorLogger.Error("error at deserializing login message", "error", err.Error())
		return nil, err
	}

	ke2, err := o.server.LoginInit(ke1, opaqueRecord)
	if err != nil {
		log.ErrorLogger.Error("error at login initializing", "error", err.Error())
		return nil, err
	}

	return ke2.Serialize(), nil
}

func (o *opaqueAdaptor) getPublicKey() ([]byte, error) {
	return os.ReadFile(o.config.Opaque.PublicKeyPath)
}

func (o *opaqueAdaptor) getPrivateKey() ([]byte, error) {
	return os.ReadFile(o.config.Opaque.PrivateKeyPath)
}

func (o *opaqueAdaptor) getOprfKey() ([]byte, error) {
	return os.ReadFile(o.config.Opaque.OprfKeyPath)
}
