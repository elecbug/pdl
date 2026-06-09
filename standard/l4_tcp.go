package standard

import "github.com/elecbug/pdl"

func TCPSource(payload pdl.PayloadFormat) pdl.Source {
	return pdl.NewSource(`
packet ` + pdl.TCP.String() + `
set mode BIG_ENDIAN MSB_FIRST

def {
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
    
    options        from 160 length (*data_offset * 32 - 160)
    payload        from (*data_offset * 32) to end
}

out json {
    src_port       tcp.source_port           DEC
    dst_port       tcp.destination_port      DEC

    seq            tcp.sequence_number       DEC
    ack            tcp.acknowledgment_number DEC

    data_offset    tcp.header_length_words   DEC
    reserved       tcp.reserved              BIN
    ns             tcp.nonce_sum             BOOL

    window         tcp.window_size           DEC
    checksum       tcp.checksum              HEX
    urgent_pointer tcp.urgent_pointer        DEC

    options        tcp.options               HEX
    payload        tcp.payload               ` + payload.String() + `

    flags<0> tcp.flag.cwr {
        0 : false
        1 : true
    }

    flags<1> tcp.flag.ece {
        0 : false
        1 : true
    }

    flags<2> tcp.flag.urg {
        0 : false
        1 : true
    }

    flags<3> tcp.flag.ack {
        0 : false
        1 : true
    }

    flags<4> tcp.flag.psh {
        0 : false
        1 : true
    }

    flags<5> tcp.flag.rst {
        0 : false
        1 : true
    }

    flags<6> tcp.flag.syn {
        0 : false
        1 : true
    }

    flags<7> tcp.flag.fin {
        0 : false
        1 : true
    }
}
`)
}
