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

## Mode

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
()
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
    data_offset from 96 length 4

    options from 160 length (*data_offset * 32 - 160)

    payload from (*data_offset * 32) to end
}
```

`*field_name` 문법은 필드의 정수값을 의미한다.

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

예시:

```pdl
out json {
    checksum checksum HEX
}
```

---

### Bit Extraction

특정 비트를 추출할 수 있다.

```pdl
out json {
    flags<6> flag.syn {
        0 : false
        1 : true
    }
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

## 현재 구현 범위

지원됨:

* packet
* set mode
* var
* 산술식 (+,-,*,/)
* def
* from
* length
* to
* end
* *field 참조
* out json
* DEC
* HEX
* BIN
* BOOL
* JSON Path
* flags<n>
* Value Mapping

---

## 향후 확장 예정

* include
* if / else
* struct
* array
* ASCII
* UTF8
* UTF16
* BASE64
* BASE58
* IP4
* IP6
* MAC
* Go Generator
* Rust Generator
* C# Generator
* Wireshark Generator
* Markdown Generator
