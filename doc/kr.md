# Packet Definition Language (PDL)

## 개요

PDL(Packet Definition Language)은 네트워크 프로토콜, 파일 포맷, 바이너리 메시지 구조를 비트 단위로 정의하기 위한 DSL(Domain Specific Language)이다.

PDL의 목표는 다음과 같다.

* 바이너리 구조 정의
* 패킷 파싱
* JSON 출력
* 프로토콜 문서 생성
* 컴파일된 스키마 저장 및 재사용

PDL은 패킷 구조를 기술하는 단일 진실 공급원(Single Source of Truth)을 목표로 한다.

---

## 설계 철학

PDL은 세 개의 계층으로 구성된다.

### 1. Variable Layer

상수 및 계산식을 정의한다.

```pdl
var {
    header_size = 160
    payload_start = header_size + 32
}
```

---

### 2. Layout Layer

비트 위치와 길이를 정의한다.

```pdl
def {
    src_port from 0 length 16
}
```

---

### 3. Presentation Layer

파싱된 값을 JSON 구조로 출력한다.

```pdl
out json {
    src_port source_port DEC
}
```

---

## 패킷 정의

모든 문서는 packet 키워드로 시작한다.

```pdl
packet TCP
```

---

## Mode 정의

PDL은 Byte Order와 Bit Order를 분리하여 정의한다.

### Byte Order

* BIG_ENDIAN
* LITTLE_ENDIAN

### Bit Order

* MSB_FIRST
* LSB_FIRST

예시:

```pdl
set mode BIG_ENDIAN MSB_FIRST
```

---

## Variable

변수는 var 블록에서 선언한다.

```pdl
var {
    a = 10
    b = a * 2
}
```

지원 연산:

```text
+
-
*
/
%
()

<<
>>
|
&
^
!

||
&&
==
!=
```

예시:

```pdl
var {
    a = 5
    b = (a + 3) * 2
}
```

---

## Layout

### Length 기반 정의

```pdl
def {
    seq from 32 length 32
}
```

---

### Range 기반 정의

```pdl
def {
    src from 0 to 15
}
```

---

### End 기반 정의

패킷 끝까지 의미한다.

```pdl
def {
    payload from 160 to end
}
```

---

### Field Value 참조

이미 정의된 필드 값을 계산식에서 사용할 수 있다.

```pdl
def {
    data_offset from 96                  length 4
    options     from 160                 length (*data_offset * 32 - 160)
    payload     from (*data_offset * 32) to     end
}
```

포인터 문법은 각 패킷의 해당 필드에 대한 실제 정수값을 의미한다.
즉, `*data_offset`은 변환하려는 패킷의 `data_offset`에 해당하는 4 bits 숫자 값이다.

---

### 분기 정의

필드의 길이를 기존에 정의된 필드에 값에 따라 서로 다르게 정의가 가능하다.

```pdl
def {
    payload switch *data_offset {
        1       : from 16 to end
        2       : from 20 to end
        3       : from 32 to end
        default : from 8  to end
    }
}
```

즉, 위 예시에서 `payload` 필드는 `data_offset`에 해당하는 값에 따라 그 길이가 달라진다.

---

## Output

출력은 out json 블록으로 정의한다.

```pdl
out json {
    src_port source_port DEC
}
```

---

### JSON Path

중첩 구조를 생성할 수 있다.

예시:

```pdl
out json {
    src_port source.port DEC
}
```

결과:

```json
{
  "source": {
    "port": 80
  }
}
```

---

### Format Type

현재 구현된 포맷은 다음과 같다.

| Format | Description        |
| ------ | ------------------ |
| DEC    | Unsigned Decimal   |
| HEX    | Hexadecimal String |
| BIN    | Binary String      |
| BOOL   | Boolean            |
| ASCII  | Pure Alphabet      |
| UTF8   | UTF8 Unicode       |
| IP4    | IP Version 4(L3)   |
| IP6    | IP Version 6(L3)   |
| MAC    | MAC Address(L2)    |

예시:

```pdl
out json {
    checksum checksum HEX
}
```

---

### Value Mapping

필드 값을 다른 값으로 매핑할 수 있다.

```pdl
out json {
    protocol protocol {
        1 : "ICMP"
        6 : "TCP"
        17 : "UDP"
        default: "Unknown"
    }
}
```

결과:

```json
{
  "protocol": "TCP"
}
```

---

### Bit Mapping

```pdl
out json {
    flags<6> flag.syn {
        0 : false
        1 : true
    }
}
```

결과:

```json
{
  "flag": {
    "syn": true
  }
}
```

---

### Type Mapping

다른 PDL을 바탕으로 출력형을 정의할 수 있다.

```pdl
out json {
    payload as TCP
}
```

결과:

```json
{
    "payload": {
        // JSON of TCP format
    }
}
```

---

### 분기 정의

출력 타입 역시 분기 정의가 가능하다.

```pdl
def {
    payload from 320 to end
}

out json {
    next_header ip.next_header {
        6       : "TCP"
        17      : "UDP"
        44      : "Fragment"
        58      : "ICMPv6"
        default : "Unknown"
    }

    payload ip.payload as switch *next_header {
        6       : TCP
        17      : UDP
        44      : IPv6Fragment
        58      : ICMPv6
        default : HEX
    }
}
```

`val` 키워드를 사용하면 조건 식으로 연산할 수도 있다.

```pdl
out json {
    src_port    udp.source_port      DEC
    dst_port    udp.destination_port DEC
    len         udp.length           DEC
    checksum    udp.checksum         HEX

    payload udp.payload as switch *dst_port {
        val == 443 || *src_port == 443 : QUIC
        default                        : HEX
    }
}
```

혹은 `as switch {...`의 형태로 selector를 지정하지 않고 `val` 부분을 사용하지 않을 수도 있다.


```pdl
out json {
    ...

    payload udp.payload as switch {
        *dst_port == 443 || *src_port == 443 : QUIC
        default                              : HEX
    }
}
```

---

### Example

```pdl
packet TCP

set mode BIG_ENDIAN MSB_FIRST

def {
    src_port    from 0 length 16
    dst_port    from 16 length 16

    data_offset from 96 length 4
    flags       from 104 length 8

    payload     from (*data_offset * 32) to end
}

out json {
    src_port source_port DEC
    dst_port destination_port DEC

    flags<6> flag.syn {
        0 : false
        1 : true
    }

    payload payload HEX
}
```

---

## 향후 확장 예정

* repeat
* bit operation
* UTF16
* BASE64
* BASE58
* Go Generator
* Rust Generator
* C# Generator
* Wireshark Generator
* Markdown Generator
