package abi

import (
	"encoding/hex"
	"fmt"
	"madledger/common/crypto/sha3"
	"madledger/common/util"
)

// Variable defines the struct of key:value
type Variable struct {
	Name  string
	Value string
}

// Packer Convenience Packing Functions
func Packer(abiData, funcName string, args ...string) ([]byte, error) {
	abiSpec, err := ReadAbiSpec([]byte(abiData))
	if err != nil {
		return nil, err
	}

	iArgs := make([]interface{}, len(args))
	for i, s := range args {
		iArgs[i] = interface{}(s)
	}
	packedBytes, err := abiSpec.Pack(funcName, iArgs...)
	if err != nil {
		return nil, err
	}

	return packedBytes, nil
}

// Unpacker read the file of abi and input the name of function and
// output of function, then unpack the variables of output
func Unpacker(abiFile, name string, data []byte) ([]*Variable, error) {
	abiSpec, err := ReadAbiSpecFile(abiFile)
	if err != nil {
		return nil, err
	}
	var args []Argument

	if name == "" {
		args = abiSpec.Constructor.Outputs
	} else {
		if _, ok := abiSpec.Functions[name]; ok {
			args = abiSpec.Functions[name].Outputs
		} else {
			args = abiSpec.Fallback.Outputs
		}
	}

	if args == nil {
		return nil, fmt.Errorf("no such function")
	}
	vars := make([]*Variable, len(args))

	if len(args) == 0 {
		return nil, nil
	}

	vals := make([]interface{}, len(args))
	for i := range vals {
		vals[i] = new(string)
	}
	err = Unpack(args, data, vals...)
	if err != nil {
		return nil, err
	}

	for i, a := range args {
		if a.Name != "" {
			vars[i] = &Variable{Name: a.Name, Value: *(vals[i].(*string))}
		} else {
			vars[i] = &Variable{Name: fmt.Sprintf("%d", i), Value: *(vals[i].(*string))}
		}
	}

	return vars, nil
}

// GetFuncHash return the hash of function
func GetFuncHash(abiFile, funcName string) (string, error) {
	abiSpec, err := ReadAbiSpecFile(abiFile)
	if err != nil {
		return "", err
	}

	if _, ok := abiSpec.Functions[funcName]; ok {
		args := abiSpec.Functions[funcName].Inputs
		var input = funcName + "("
		for _, a := range args {
			input += a.EVM.GetSignature()
		}
		input += ")"
		hash := sha3.Sha3([]byte(input))
		return util.Hex(hash[:4]), nil
	}
	return "", fmt.Errorf("no such function")
}

// GetPayload return the payload string
func GetPayload(abiFile, funcName string, inputs []string) (string, error) {
	abiSpec, err := ReadAbiSpecFile(abiFile)
	if err != nil {
		return "", err
	}

	if _, ok := abiSpec.Functions[funcName]; ok {
		args := abiSpec.Functions[funcName].Inputs
		if len(args) != len(inputs) {
			return "", fmt.Errorf("Except %d inputs other than %d inputs", len(args), len(inputs))
		}
		var input = funcName + "("
		var payload []byte
		for i, a := range args {
			input += a.EVM.GetSignature()
			bs, err := a.EVM.pack(inputs[i])
			if err != nil {
				return "", err
			}
			payload = util.BytesCombine(payload, bs)
		}
		input += ")"
		hash := sha3.Sha3([]byte(input))
		return util.Hex(hash[:4]) + util.Hex(payload), nil
	}
	return "", fmt.Errorf("no such function")
}

// GetPayloadBytes return the payload bytes
func GetPayloadBytes(abiFile, funcName string, inputs []string) ([]byte, error) {
	payload, err := GetPayload(abiFile, funcName, inputs)
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(payload)
}

func stripHex(s string) string {
	if len(s) > 1 {
		if s[:2] == "0x" {
			s = s[2:]
			if len(s)%2 != 0 {
				s = "0" + s
			}
			return s
		}
	}
	return s
}

// func ReadAbiFormulateCallFile(abiLocation, binPath, funcName string, args []string) ([]byte, error) {
// 	abiSpecBytes, err := readAbi(binPath, abiLocation)
// 	if err != nil {
// 		return []byte{}, err
// 	}

// 	return Packer(abiSpecBytes, funcName, args...)
// }

// func ReadAbiFormulateCall(abiSpecBytes []byte, funcName string, args []string) ([]byte, error) {
// 	return Packer(string(abiSpecBytes), funcName, args...)
// }

// func ReadAndDecodeContractReturn(abiLocation, binPath, funcName string, resultRaw []byte) ([]*Variable, error) {
// 	abiSpecBytes, err := readAbi(binPath, abiLocation)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Unpack the result
// 	return Unpacker(abiSpecBytes, funcName, resultRaw)
// }

// func readAbi(root, contract string) (string, error) {
// 	p := path.Join(root, stripHex(contract))
// 	if _, err := os.Stat(p); err != nil {
// 		p = path.Join(root, stripHex(contract)+".bin")
// 		if _, err = os.Stat(p); err != nil {
// 			return "", fmt.Errorf("Abi doesn't exist for =>\t%s", p)
// 		}
// 	}

// 	b, err := ioutil.ReadFile(p)
// 	if err != nil {
// 		return "", err
// 	}
// 	sol := SolidityOutputContract{}
// 	err = json.Unmarshal(b, &sol)
// 	if err != nil {
// 		return "", err
// 	}

// 	return string(sol.Abi), nil
// }

// func Unpacker(abiData, name string, data []byte) ([]*Variable, error) {
// 	abiSpec, err := ReadAbiSpec([]byte(abiData))
// 	if err != nil {
// 		return nil, err
// 	}

// 	var args []Argument

// 	if name == "" {
// 		args = abiSpec.Constructor.Outputs
// 	} else {
// 		if _, ok := abiSpec.Functions[name]; ok {
// 			args = abiSpec.Functions[name].Outputs
// 		} else {
// 			args = abiSpec.Fallback.Outputs
// 		}
// 	}

// 	if args == nil {
// 		return nil, fmt.Errorf("no such function")
// 	}
// 	vars := make([]*Variable, len(args))

// 	if len(args) == 0 {
// 		return nil, nil
// 	}

// 	vals := make([]interface{}, len(args))
// 	for i, _ := range vals {
// 		vals[i] = new(string)
// 	}
// 	err = Unpack(args, data, vals...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for i, a := range args {
// 		if a.Name != "" {
// 			vars[i] = &Variable{Name: a.Name, Value: *(vals[i].(*string))}
// 		} else {
// 			vars[i] = &Variable{Name: fmt.Sprintf("%d", i), Value: *(vals[i].(*string))}
// 		}
// 	}

// 	return vars, nil
// }
