# Packet Definition Language (PDL)

## 개요

PDL(Packet Definition Language)은 네트워크 프로토콜, 파일 포맷, 바이너리 메시지 구조를 비트 단위로 정의하기 위한 DSL(Domain Specific Language)이다.

PDL의 목표는 다음과 같다.

* 바이너리 구조 정의
* 패킷 파싱
* JSON 출력
* 프로토콜 문서 생성
* Wireshark Dissector 생성
* Go / Rust / C# 코드 생성

PDL은 "패킷 구조를 기술하는 단일 진실 공급원(Single Source of Truth)"을 목표로 한다.

---

## 설계 철학

PDL은 세 가지 계층으로 분리된다.

### 1. Define Layer

기본 변수 및 값 정의

예시

```pdl
var {
    a = 50
}
```

### 2. Layout Layer

바이트 및 비트 위치 정의

예시

```pdl
def {
    src_port from 0 length 16
}
```

### 3. Presentation Layer

JSON 출력 방식 정의

예시

```pdl
out json {
    src_port source_port DEC
}
```

---

## 문법

### 주석

주석은 `#`을 이용한다.

---

### 패킷

모든 프로토콜은 `packet`과 해당 패킷의 이름으로 시작한다.

예시: UDP 패킷 정의

```pdl
packet UDP

var {...}
def {...}
out json {...}
```

---

### 데이터

데이터 정의는 `var` 블럭을 사용한다.

예시: 변수 `a`, `b`, `c` 등을 선언하고 값을 할당

```pdl
var {
    a = 50
    b = a * 2  # b = 100
    c = a + b  # c = 150
}
```

---

### 정의

정의는 `def` 블럭을 통해 이루어지며, 블럭 내 각 한줄마다 의미를 가진다.
각 줄은 변수 이름과 `from`, `to`, `length` 키워드 등을 통해 구성된다.

#### 길이 기반 정의

예시: 변수 `seq`를 32번 비트부터 시작하여 32개의 비트로 정의 (즉, 전체 비트 중 [32] [33] [34] ... [62] [63])

```pdl
def {
    seq from 32 length 32  # or `seq from 32;32`
}
```

#### 범위 기반 정의

예시: 변수 `src`를 0번 비트부터 시작하여 15번 비트까지로 정의 (즉, 전체 비트 중 [0] [1] [2] ... [14] [15])

```pdl
def {
    src from 0 to 15  # or `src from 0:15`
}
```

혹은 `to end`를 사용하여 남은 패킷 전체를 의미할 수도 있다.

예시: 변수 `payload`를 200번 비트부터 마지막까지로 정의

```pdl
def {
    payload from 200 to end
}
```

#### 다른 변수 기반 정의

다른 변수의 값을 활용하여 변수의 길이를 정할 수도 있으며, 다음과 같은 구조를 따른다.

예시: `data_offset`에 할당되는 값을 Little Endian 방식의 정수로 해석하여, 그 값을 바탕으로 `options` 변수의 범위를 지정

```pdl
def {
    data_offset from 96  length 4
    options     from 160 length (*data_offset * 32 - 160) as LITTLE_ENDIAN
}
```

---

### 출력

출력은 `out` 블럭을 통해 이루어지며, 블럭 내 각 한줄마다 의미를 가진다.
`out`은 [정의](#정의)에서 정의한 각 필드를 JSON 등과 어떻게 대응 시킬지 정의한다.
각 줄은 변수 이름과 JSON 이름, Format Type 키워드 등을 통해 구성되며, 의미는 다음과 같다.

예시: 변수 `seq`와 JSON `"sequence_number"` 필드를 대응시키고, 값은 10진수 숫자 포맷으로 해석 및 대입하여 `"sequence_number": 11`의 형태로 출력

```pdl
out json {
    seq sequence_number DEC
}
```

JSON 이름 부분을 조정하여 JSON 구조를 설계할 수도 있다.

* `seq sequence.number DEC` => `{"sequence":{"number":11}}`
* `seq sequence[0].number DEC` => `{"sequence":[{"number":11}]}`
* `seq [0] DEC` => `[11]`

#### Endian

PDL은 패킷 단위 Endian을 지원한다.

각 `out` 블록은 자신의 위 가장 가까운 `set mode`의 Endian 방식을 따르며, `BIG_ENDIAN` 또는 `LITTLE_ENDIAN`으로 정의된다.

만약 `value`가 `0b11001010`이라면, `BIG_ENDIAN` 기준

* `value<0>` -> 1
* `value<1>` -> 1
* `value<2>` -> 0
* `value<3>` -> 0
* ...
* `value<6>` -> 1
* `value<7>` -> 0

예시

```pdl
packet ExampleProtocol

def {...}

set mode BIG_ENDIAN     # Can use `set mode LITTLE_ENDIAN`
out json {...}
```

#### Format Type (e.g. Based `LITTLE_ENDIAN`)

|Format Type|Description                  |Example Bits         |JSON Output     |
|-----------|-----------------------------|---------------------|----------------|
|`DEC`      |Unsigned Decimal Integer     |`00001110`           |`14`            |
|`HEX`      |Unsigned Hexadecimal Integer |`00001110`           |`"0E"`          |
|`BIN`      |Unsigned Binary Integer      |`00001110`           |`"00001110"`    |
|`BOOL`     |Boolean Type (Only 1 Bit)    |`0`, `1`             |`false`, `true` |
|`ASCII`    |ASCII Code Encoding          |`0110001101100001`   |`"ca"`          |
|`UTF8`     |UTF8 Encoding String         |
|`UTF16`    |UTF16 Encoding String        |
|`BASE64`   |Base64 Encoded String        |
|`BASE58`   |Base58 Encoded String        |
|`IP4`      |IPv4 Type Bytes              |
|`IP6`      |IPv6 Type Bytes              |
|`MAC`      |MAC Address Type Bytes       |

#### Bracket Out

`out`은 값 매핑을 통해 출력 형태를 결정할 수도 있다.

단일 비트 Flag 기반 출력 예시

```pdl
out json {
    flags<6> syn {
        0 : false
        1 : true
    }
}
```

다중 비트 Flag 기반 출력 예시

```pdl
out json {
    l3_protocol_type protocol {
        1: 
        2:
        3:
        4:
    }
}
```

---

## [TCP 정의 예시](./eg_tcp.pdl)

---

## 향후 확장 예정

### 조건부 필드

```pdl
if protocol == 1 {
    payload from 0 to 16
}
else if protocol == 2 {
    payload from 0 to 32
}
else {
    payload from 0 to 8
}
```

---

## 코드 작성 구조

```text
PDL
 ↓
Parser
 ↓
AST
 ↓
 ├─ JSON Decoder
 ├─ Go Generator
 ├─ Rust Generator
 ├─ C# Generator
 ├─ Wireshark Lua Generator
 ├─ Markdown Documentation Generator
 └─ Packet Visualizer
```
