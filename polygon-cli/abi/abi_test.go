package abi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHashFunctionSelector tests function signature from various common signatures documented by openchain.xyz:
// https://openchain.xyz/signatures
// Future tests:
// - pass function types that doesn't exist
func TestHashFunctionSelector(t *testing.T) {
	type Test struct {
		Name              string
		FunctionSignature string
		Expected          string
		ErrMsg            string
	}

	tests := []Test{
		{
			Name:              "send",
			FunctionSignature: "send(address,uint64,address,uint256,uint64,bytes32,string)",
			Expected:          "037d684b",
		},
		{
			Name:              "sendTxns",
			FunctionSignature: "sendTxns(address,(uint256,address,address),(bytes,bytes),(string,address,uint256,bytes)[])",
			Expected:          "63486689",
		},
		{
			Name:              "mint",
			FunctionSignature: "minted(uint256,string,string,string,uint256,string,string,string,string[],(string,string))",
			Expected:          "019015a3",
		},
		{
			Name:              "mintSuccessful",
			FunctionSignature: "mintSuccessful(address,uint256,uint256,bytes)",
			Expected:          "001d98a3",
		},
		{
			Name:              "sendForTokens",
			FunctionSignature: "sendForTokens(uint256,address[])",
			Expected:          "1114fd36",
		},
		{
			Name:              "sendTransfer",
			FunctionSignature: "sendTransfer(address,(uint256,address,address),(bytes,bytes),(address,address,uint256,uint256))",
			Expected:          "4b776c6d",
		},
		{
			Name:              "Invalid Parenthesis Error",
			FunctionSignature: "sendTransfer(address,uint256",
			ErrMsg:            "invalid parenthesis",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			output, err := HashFunctionSelector(tc.FunctionSignature)
			if tc.ErrMsg != "" {
				assert.ErrorContains(t, err, tc.ErrMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, output, tc.Expected)
			}
		})
	}
}

// TestAbiEncode tests various encoding of function signatures and their inputs.
// The expected values are derived from `cast calldata`.
// Future tests:
// - int value outside of int bounds
// - uint value outside of uint bounds
// - uint value outside of uint bounds
// - For a fixed sized array, type[M], pass invalid length
// - pass function types that doesn't exist
func TestAbiEncode(t *testing.T) {
	type AbiEncodeInput struct {
		FunctionSignature string
		FunctionInputs    []string
	}

	type Test struct {
		Name     string
		Input    AbiEncodeInput
		Expected string
		ErrMsg   string
	}

	tests := []Test{
		{
			// cast calldata "f(uint8)" 19
			Name: "Simple uint8 Encoding",
			Input: AbiEncodeInput{
				FunctionSignature: "f(uint8)",
				FunctionInputs: []string{
					"19",
				},
			},
			Expected: "0x3120d4340000000000000000000000000000000000000000000000000000000000000013",
		},
		{
			// cast calldata "f(uint256)" 1999
			Name: "Simple uint256 Encoding",
			Input: AbiEncodeInput{
				FunctionSignature: "f(uint256)",
				FunctionInputs: []string{
					"1999",
				},
			},
			Expected: "0xb3de648b00000000000000000000000000000000000000000000000000000000000007cf",
		},
		{
			Name: "Error: Negative uint256 Encoding",
			Input: AbiEncodeInput{
				FunctionSignature: "f(uint256)",
				FunctionInputs: []string{
					"-1999",
				},
			},
			ErrMsg: "can't be negative",
		},
		{
			// cast calldata "f(int8, int32, int256, int256)" 99 999 999999 -999999
			Name: "int Encoding with positive and negative num",
			Input: AbiEncodeInput{
				FunctionSignature: "f(int8, int32, int256, int256)",
				FunctionInputs: []string{
					"99",
					"999",
					"999999",
					"-999999",
				},
			},
			Expected: "0x15842b5c000000000000000000000000000000000000000000000000000000000000006300000000000000000000000000000000000000000000000000000000000003e700000000000000000000000000000000000000000000000000000000000f423ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0bdc1",
		},
		{
			// cast calldata "f(bool,bool)" true false
			Name: "bool Encoding",
			Input: AbiEncodeInput{
				FunctionSignature: "f(bool,bool)",
				FunctionInputs: []string{
					"true",
					"false",
				},
			},
			Expected: "0xad51369a00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			Name: "Error: bool Encoding with invalid input",
			Input: AbiEncodeInput{
				FunctionSignature: "f(bool,bool)",
				FunctionInputs: []string{
					"true",
					"no",
				},
			},
			ErrMsg: "bool must be either 'true' or 'false'",
		},
		{
			// cast calldata "f(address,address)" 0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6 85dA99c8a7C2C95964c8EfD687E95E632Fc533D6
			Name: "address Encoding",
			Input: AbiEncodeInput{
				FunctionSignature: "f(address, address)",
				FunctionInputs: []string{
					"0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6",
					"85dA99c8a7C2C95964c8EfD687E95E632Fc533D6",
				},
			},
			Expected: "0x4d201ccb00000000000000000000000085da99c8a7c2c95964c8efd687e95e632fc533d600000000000000000000000085da99c8a7c2c95964c8efd687e95e632fc533d6",
		},
		{
			// cast calldata "f(bytes3,bytes5,bytes)" 0x123456 1234567890 ffffffff88888888888ffff111
			Name: "bytes Encoding",
			Input: AbiEncodeInput{
				FunctionSignature: "f(bytes3,bytes5,bytes)",
				FunctionInputs: []string{
					"0x123456",
					"1234567890",
					"ffffffff88888888888ffff111",
				},
			},
			Expected: "0x4f0f2614123456000000000000000000000000000000000000000000000000000000000012345678900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000000dffffffff88888888888ffff11100000000000000000000000000000000000000",
		},
		{
			Name: "Error: invalid bytes<M> Encoding",
			Input: AbiEncodeInput{
				FunctionSignature: "f(bytes3)",
				FunctionInputs: []string{
					"0x1234567",
				},
			},
			ErrMsg: "Invalid string length",
		},
		{
			Name: "Error: invalid bytes Encoding. Odd number of digits",
			Input: AbiEncodeInput{
				FunctionSignature: "f(bytes)",
				FunctionInputs: []string{
					"0x1234567",
				},
			},
			ErrMsg: "Odd number of digits",
		},
		{
			// cast calldata "f(string)(string)" "adfjkadhsffdhjksfdahjsfhadjsfasdhjfdsjlkfadshkjladfshjkadfskjladsfjkldfajhkdjafhkadsfjkldjksafjkhldsfhjksadflhj kldsafjklhadfsjkahlsdfkjlhasdfjkadfhslajkhsadfsjkl"
			Name: "Simple String Encoding",
			Input: AbiEncodeInput{
				FunctionSignature: "f(string)(string)",
				FunctionInputs: []string{
					"adfjkadhsffdhjksfdahjsfhadjsfasdhjfdsjlkfadshkjladfshjkadfskjladsfjkldfajhkdjafhkadsfjkldjksafjkhldsfhjksadflhj kldsafjklhadfsjkahlsdfkjlhasdfjkadfhslajkhsadfsjkl",
				},
			},
			Expected: "0x91e145ef000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000a26164666a6b61646873666664686a6b73666461686a73666861646a7366617364686a6664736a6c6b66616473686b6a6c61646673686a6b616466736b6a6c616473666a6b6c6466616a686b646a6166686b616473666a6b6c646a6b7361666a6b686c647366686a6b736164666c686a206b6c647361666a6b6c68616466736a6b61686c7364666b6a6c68617364666a6b61646668736c616a6b6873616466736a6b6c000000000000000000000000000000000000000000000000000000000000",
		},
		{
			// cast calldata "f(string)" ""
			Name: "Empty String",
			Input: AbiEncodeInput{
				FunctionSignature: "f(string)(string)",
				FunctionInputs: []string{
					"",
				},
			},
			Expected: "0x91e145ef00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			// cast calldata "f(string[])" '[]'
			Name: "Empty String Array Encoding",
			Input: AbiEncodeInput{
				FunctionSignature: "f(string[])(string)",
				FunctionInputs: []string{
					"[]",
				},
			},
			Expected: "0xe9cc878000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			// cast calldata "f(uint256[])" '[12,34,567]'
			Name: "Simple Array Encoding",
			Input: AbiEncodeInput{
				FunctionSignature: "f(uint256[])(int8)",
				FunctionInputs: []string{
					"[12,34,567]",
				},
			},
			Expected: "0x7bc5bbbf00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000220000000000000000000000000000000000000000000000000000000000000237",
		},
		{
			Name: "Simple Fixed Array Size Encoding",
			Input: AbiEncodeInput{
				FunctionSignature: "f(uint256[3])(int8)",
				FunctionInputs: []string{
					"[12,34,567]",
				},
			},
			Expected: "0x5dc128910000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000220000000000000000000000000000000000000000000000000000000000000237",
		},
		{
			Name: "Invalid Array Input: Unclosed bracket",
			Input: AbiEncodeInput{
				FunctionSignature: "f(uint256[])(int8)",
				FunctionInputs: []string{
					"[12,34,567",
				},
			},
			ErrMsg: `expected "]"`,
		},
		{
			Name: "Invalid Array Input: Mismatched type",
			Input: AbiEncodeInput{
				FunctionSignature: "f(uint256[])(int8)",
				FunctionInputs: []string{
					`[12,34,"yes"]`,
				},
			},
			ErrMsg: `Failed to convert`,
		},
		{
			// cast calldata "f(uint256[][])" '[[12,34,567],[987,654,321,0],[99999999,99999]]'
			Name: "Multidimensional Array Encoding",
			Input: AbiEncodeInput{
				FunctionSignature: "f(uint256[][])",
				FunctionInputs: []string{
					"[[12,34,567],[987,654,321,0],[99999999,99999]]",
				},
			},
			Expected: "0xc26b6b9a00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000001800000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000220000000000000000000000000000000000000000000000000000000000000237000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000003db000000000000000000000000000000000000000000000000000000000000028e0000000000000000000000000000000000000000000000000000000000000141000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000005f5e0ff000000000000000000000000000000000000000000000000000000000001869f",
		},
		{
			// cast calldata "newSchool((string,string[],uint256,bool))" '("matic",["123 street ave.","321 ave st."], 9999, false)'
			Name: "Struct Encoding",
			Input: AbiEncodeInput{
				FunctionSignature: "newSchool((string,string[],uint256,bool))",
				FunctionInputs: []string{
					`("matic",["123 street ave.","321 ave st."], 9999, false)`,
				},
			},
			Expected: "0x5866fb060000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000270f000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000056d61746963000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000f31323320737472656574206176652e0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000b333231206176652073742e000000000000000000000000000000000000000000",
		},
		{
			Name: "Invalid Struct Input: Unclosed Tuple",
			Input: AbiEncodeInput{
				FunctionSignature: "newSchool((string,string[],uint256,bool))",
				FunctionInputs: []string{
					`("matic",["123 street ave.","321 ave st."], 9999, false`,
				},
			},
			ErrMsg: `expected ")"`,
		},
		{
			Name: "Invalid Struct Input: invalid num of elements",
			Input: AbiEncodeInput{
				FunctionSignature: "newSchool((string,string[],uint256,bool))",
				FunctionInputs: []string{
					`("matic",["123 street ave.","321 ave st."], 9999)`,
				},
			},
			ErrMsg: `Mismatched length of tuple elements`,
		},
		{
			// cast calldata "getMultipliedAndAddNumber(uint256,bytes3,bool,string,address,int256[],string[],string[][],(string,uint256,bool[]))" 100000 "0x123456" true "abc" "0x6fda56c57b0acadb96ed5624ac500c0429d59429" "[1,2,3,4]" '["hi","bye","YOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOZ","test"]' '[["hi","bye","test"], ["hhhhh","byebye","bytebye","bytebyte"]]' '("yesyesyesyes", 369, [true,false,true])'
			Name: "A complex encoding",
			Input: AbiEncodeInput{
				FunctionSignature: "getMultipliedAndAddNumber(uint256,bytes3,bool,string,address,int256[],string[],string[][],(string,uint256,bool[]))",
				FunctionInputs: []string{
					"100000",
					"0x123456",
					"true",
					"abc",
					"0x6fda56c57b0acadb96ed5624ac500c0429d59429",
					"[1,2,3,4]",
					`["hi","bye","YOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOZ","test"]`,
					`[["hi","bye","test"], ["hhhhh","byebye","bytebye","bytebyte"]]`,
					`("yesyesyesyes", 369, [true,false,true])`,
				},
			},
			Expected: `0xaf1174e400000000000000000000000000000000000000000000000000000000000186a01234560000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000001200000000000000000000000006fda56c57b0acadb96ed5624ac500c0429d594290000000000000000000000000000000000000000000000000000000000000160000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000003c0000000000000000000000000000000000000000000000000000000000000070000000000000000000000000000000000000000000000000000000000000000036162630000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000016000000000000000000000000000000000000000000000000000000000000000026869000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000362796500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000022594f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f5a000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000474657374000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000001800000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000e00000000000000000000000000000000000000000000000000000000000000002686900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000036279650000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000474657374000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000000000000000000000000000000568686868680000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006627965627965000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000076279746562796500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000862797465627974650000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000017100000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000000c79657379657379657379657300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			output, err := AbiEncode(tc.Input.FunctionSignature, tc.Input.FunctionInputs)
			if tc.ErrMsg != "" {
				assert.ErrorContains(t, err, tc.ErrMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, output, tc.Expected)
			}
		})
	}
}
