package standard

import "github.com/elecbug/pdl"

func DNSSource(payload pdl.PayloadFormat) pdl.Source {
	_ = payload

	return pdl.NewSource(`
packet ` + pdl.DNS.String() + `

set mode BIG_ENDIAN MSB_FIRST

def {
    transaction_id from 0  length 16
    flags          from 16 length 16

    qr             from 16 length 1
    opcode         from 17 length 4
    aa             from 21 length 1
    tc             from 22 length 1
    rd             from 23 length 1
    ra             from 24 length 1
    z              from 25 length 3
    rcode          from 28 length 4

    qdcount        from 32 length 16
    ancount        from 48 length 16
    nscount        from 64 length 16
    arcount        from 80 length 16

    payload        from 96 to end
}

out json {
    transaction_id dns.transaction_id HEX

    flags          dns.flags HEX

    qr dns.is_response {
        0 : false
        1 : true
    }

    opcode dns.opcode {
        0       : "Query"
        1       : "IQuery"
        2       : "Status"
        4       : "Notify"
        5       : "Update"
        default : "Unknown"
    }

    aa dns.authoritative_answer BOOL
    tc dns.truncated BOOL
    rd dns.recursion_desired BOOL
    ra dns.recursion_available BOOL

    z dns.reserved BIN

    rcode dns.response_code {
        0       : "NoError"
        1       : "FormErr"
        2       : "ServFail"
        3       : "NXDomain"
        4       : "NotImp"
        5       : "Refused"
        default : "Unknown"
    }

    qdcount dns.question_count DEC
    ancount dns.answer_count DEC
    nscount dns.authority_count DEC
    arcount dns.additional_count DEC

    payload dns.payload HEX
}
`)
}
