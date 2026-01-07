package opaque

import (
	"os"

	"github.com/TheAmirhosssein/cool-password-manage/config"
	"github.com/bytemare/opaque"
)

type OpaqueAdaptor struct {
	server *opaque.Server
	config *config.Config
}

func New(config *config.Config) (OpaqueService, error) {
	a := &OpaqueAdaptor{config: config}
	return a, a.Init()
}

func (o *OpaqueAdaptor) Init() error {
	serverID := []byte(o.config.Opaque.ServerID)

	conf := opaque.DefaultConfiguration()

	secretOprfSeed, err := o.getOprfKey()
	if err != nil {
		return err
	}

	serverPrivateKey, err := o.getPrivateKey()
	if err != nil {
		return err
	}

	serverPublicKey, err := o.getPublicKey()
	if err != nil {
		return err
	}

	server, err := conf.Server()
	if err != nil {
		return err
	}

	if err := server.SetKeyMaterial(serverID, serverPrivateKey, serverPublicKey, secretOprfSeed); err != nil {
		return err
	}

	o.server = server
	return nil
}

func (o *OpaqueAdaptor) RegisterInit(message []byte) ([]byte, []byte, error) {
	req, err := o.server.Deserialize.RegistrationRequest(message)
	if err != nil {
		return nil, nil, err
	}

	credID := opaque.RandomBytes(64)

	serverPublicKey, err := o.getPublicKey()
	if err != nil {
		return nil, nil, err
	}

	pks, err := o.server.Deserialize.DecodeAkePublicKey(serverPublicKey)
	if err != nil {
		return nil, nil, err
	}

	secretOprfKey, err := o.getOprfKey()
	if err != nil {
		return nil, nil, err
	}

	resp := o.server.RegistrationResponse(req, pks, credID, secretOprfKey)

	return resp.Serialize(), credID, nil
}

func (o *OpaqueAdaptor) RegisterFinalize(message, credID []byte, username string) ([]byte, error) {
	record, err := o.server.Deserialize.RegistrationRecord(message)
	if err != nil {
		return nil, err
	}

	clientRecord := &opaque.ClientRecord{
		CredentialIdentifier: credID,           // used during serialization
		ClientIdentity:       []byte(username), // used during serialization
		RegistrationRecord:   record,
	}

	return clientRecord.Serialize(), nil
}

func (o *OpaqueAdaptor) getPublicKey() ([]byte, error) {
	return os.ReadFile(o.config.Opaque.PublicKeyPath)
}

func (o *OpaqueAdaptor) getPrivateKey() ([]byte, error) {
	return os.ReadFile(o.config.Opaque.PrivateKeyPath)
}

func (o *OpaqueAdaptor) getOprfKey() ([]byte, error) {
	return os.ReadFile(o.config.Opaque.OprfKeyPath)
}
