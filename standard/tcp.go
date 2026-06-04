package standard

import "github.com/elecbug/pdl"

func TCPPDL(payload pdl.Payload) pdl.Source {
	return pdl.Source(`
packet ` + TCP.String() + `
set mode BIG_ENDIAN MSB_FIRST

def {
    # Fixed header: 20 bytes
    src_port       from 0   length 16
    dst_port       from 16  length 16
    
    seq            from 32  length 32
    ack            from 64  length 32
    
    data_offset    from 96  length 4
    reserved       from 100 length 3
    ns             from 103 length 1
    flags          from 104 length 8
    
    window         from 112 length 16
    checksum       from 128 length 16
    urgent_pointer from 144 length 16
    
    # Variable area
    options        from 160 length ((*data_offset * 32) - 160)
    payload        from (*data_offset * 32) to end
}

out json {
    src_port       source_port           DEC
    dst_port       destination_port      DEC

    seq            sequence_number       DEC
    ack            acknowledgment_number DEC

    data_offset    header_length_words   DEC
    reserved       reserved              BIN
    ns             nonce_sum             BOOL

    window         window_size           DEC
    checksum       checksum              HEX
    urgent_pointer urgent_pointer        DEC

    options        options               HEX
    payload        payload               ` + payload.String() + `

    # In BIG_ENDIAN mode, flags<0> is MSB and flags<7> is LSB.
    flags<0> flag.cwr {
        0 : false
        1 : true
    }

    flags<1> flag.ece {
        0 : false
        1 : true
    }

    flags<2> flag.urg {
        0 : false
        1 : true
    }

    flags<3> flag.ack {
        0 : false
        1 : true
    }

    flags<4> flag.psh {
        0 : false
        1 : true
    }

    flags<5> flag.rst {
        0 : false
        1 : true
    }

    flags<6> flag.syn {
        0 : false
        1 : true
    }

    flags<7> flag.fin {
        0 : false
        1 : true
    }
}
`)
}
