// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package hclsyntax

import (
	"reflect"
	"testing"

	"github.com/hashicorp/hcl/v2"
)

func init() {
	// see the documentation of this variable for more information
	tracePeekerNewlinesStack = true
}

func TestPeeker(t *testing.T) {
	tokens := Tokens{
		{
			Type: TokenIdent,
		},
		{
			Type: TokenComment,
		},
		{
			Type: TokenIdent,
		},
		{
			Type: TokenComment,
		},
		{
			Type: TokenIdent,
		},
		{
			Type: TokenNewline,
		},
		{
			Type: TokenIdent,
		},
		{
			Type: TokenNewline,
		},
		{
			Type: TokenIdent,
		},
		{
			Type: TokenNewline,
		},
		{
			Type: TokenEOF,
		},
	}

	{
		peeker := newPeeker(tokens, true)

		wantTypes := []TokenType{
			TokenIdent,
			TokenComment,
			TokenIdent,
			TokenComment,
			TokenIdent,
			TokenNewline,
			TokenIdent,
			TokenNewline,
			TokenIdent,
			TokenNewline,
			TokenEOF,
		}
		var gotTypes []TokenType

		for {
			peeked := peeker.Peek()
			read := peeker.Read()
			if peeked.Type != read.Type {
				t.Errorf("mismatched Peek %s and Read %s", peeked, read)
			}

			gotTypes = append(gotTypes, read.Type)

			if read.Type == TokenEOF {
				break
			}
		}

		if !reflect.DeepEqual(gotTypes, wantTypes) {
			t.Errorf("wrong types\ngot:  %#v\nwant: %#v", gotTypes, wantTypes)
		}
	}

	{
		peeker := newPeeker(tokens, false)

		wantTypes := []TokenType{
			TokenIdent,
			TokenIdent,
			TokenIdent,
			TokenNewline,
			TokenIdent,
			TokenNewline,
			TokenIdent,
			TokenNewline,
			TokenEOF,
		}
		var gotTypes []TokenType

		for {
			peeked := peeker.Peek()
			read := peeker.Read()
			if peeked.Type != read.Type {
				t.Errorf("mismatched Peek %s and Read %s", peeked, read)
			}

			gotTypes = append(gotTypes, read.Type)

			if read.Type == TokenEOF {
				break
			}
		}

		if !reflect.DeepEqual(gotTypes, wantTypes) {
			t.Errorf("wrong types\ngot:  %#v\nwant: %#v", gotTypes, wantTypes)
		}
	}

	{
		peeker := newPeeker(tokens, false)

		peeker.PushIncludeNewlines(false)

		wantTypes := []TokenType{
			TokenIdent,
			TokenIdent,
			TokenIdent,
			TokenIdent,
			TokenIdent,
			TokenNewline, // we'll pop off the PushIncludeNewlines before we get here
			TokenEOF,
		}
		var gotTypes []TokenType

		idx := 0
		for {
			peeked := peeker.Peek()
			read := peeker.Read()
			if peeked.Type != read.Type {
				t.Errorf("mismatched Peek %s and Read %s", peeked, read)
			}

			gotTypes = append(gotTypes, read.Type)

			if read.Type == TokenEOF {
				break
			}

			if idx == 4 {
				peeker.PopIncludeNewlines()
			}

			idx++
		}

		if !reflect.DeepEqual(gotTypes, wantTypes) {
			t.Errorf("wrong types\ngot:  %#v\nwant: %#v", gotTypes, wantTypes)
		}
	}
}

func TestPeeker_PrevRange(t *testing.T) {
	tokens := Tokens{
		{
			Type: TokenComment,
			Range: hcl.Range{
				Filename: "test.hcl",
				Start:    hcl.Pos{Line: 1, Column: 1, Byte: 0},
				End:      hcl.Pos{Line: 3, Column: 4, Byte: 20},
			},
		},
		{
			Type: TokenIdent,
			Range: hcl.Range{
				Filename: "test.hcl",
				Start:    hcl.Pos{Line: 4, Column: 1, Byte: 21},
				End:      hcl.Pos{Line: 4, Column: 4, Byte: 24},
			},
		},
		{
			Type: TokenEOF,
			Range: hcl.Range{
				Filename: "test.hcl",
				Start:    hcl.Pos{Line: 4, Column: 4, Byte: 24},
				End:      hcl.Pos{Line: 4, Column: 4, Byte: 24},
			},
		},
	}

	p := newPeeker(tokens, false)

	// When p.NextIndex == 0, PrevRange should anchor to the start of the first token.
	got := p.PrevRange()
	want := hcl.Range{
		Filename: "test.hcl",
		Start:    hcl.Pos{Line: 1, Column: 1, Byte: 0},
		End:      hcl.Pos{Line: 1, Column: 1, Byte: 0},
	}
	if got != want {
		t.Fatalf("wrong PrevRange at start: got %#v, want %#v", got, want)
	}

	// Reading the first non-comment token advances NextIndex.
	tok := p.Read()
	if tok.Type != TokenIdent {
		t.Fatalf("expected TokenIdent, got %s", tok.Type)
	}
	got = p.PrevRange()
	want = tokens[1].Range
	if got != want {
		t.Fatalf("wrong PrevRange after reading token: got %#v, want %#v", got, want)
	}
}

