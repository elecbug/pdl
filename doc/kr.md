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

## 실행 모델

PDL 문서는 다음 순서로 처리된다.

1. `packet` 선언을 읽는다.
2. `set mode`를 통해 기본 해석 모드를 설정한다.
3. `var` 블록의 상수 및 계산 변수를 평가한다.
4. `def` 블록을 통해 바이너리 입력에서 필드를 추출한다.
5. `out` 블록을 통해 추출된 필드를 JSON 등의 출력 구조로 변환한다.

단, `def` 블록 내부에서는 앞에서 정의된 필드의 값을 `*field_name` 형태로 참조할 수 있다.

예시

```pdl
def {
    data_offset from 96 length 4
    options from 160 length (*data_offset * 32 - 160)
}
```

위 예시에서 `*data_offset`은 `data_offset` 필드 자체가 아니라, 입력 패킷에서 추출된 `data_offset`의 값을 의미한다.

---

## 문법

### 주석

주석은 `#`을 이용한다.

```pdl
# This is a comment
src_port from 0 length 16  # Inline comment
```

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

### 모드

`set mode`는 이후 등장하는 `def`, `var`, `out` 블록의 기본 해석 방식을 지정한다.

```pdl
set mode BIG_ENDIAN
```

또는

```pdl
set mode LITTLE_ENDIAN
```

권장 구조는 패킷 선언 직후 기본 모드를 지정하는 것이다.

```pdl
packet TCP

set mode BIG_ENDIAN

def {...}
out json {...}
```

`set mode`는 다시 선언될 수 있으며, 이후 블록은 가장 가까운 위쪽 `set mode`를 따른다.

예시

```pdl
set mode BIG_ENDIAN
out json {
    value number DEC
}

set mode LITTLE_ENDIAN
out json {
    value reversed_number DEC
}
```

---

### 데이터

데이터 정의는 `var` 블록을 사용한다.

예시: 변수 `a`, `b`, `c` 등을 선언하고 값을 할당

```pdl
var {
    a = 50
    b = a * 2  # b = 100
    c = a + b  # c = 150
}
```

`var`는 주로 상수, 고정 오프셋, 반복적으로 사용되는 계산식을 정의하는 데 사용한다.

예시

```pdl
var {
    tcp_fixed_header_bits = 160
}
```

---

### 필드 값 참조

`def`에서 정의된 필드의 실제 값을 참조할 때는 `*`를 사용한다.

```pdl
*data_offset
```

이는 `data_offset`이라는 필드의 위치 정보가 아니라, 입력 패킷에서 추출된 값을 의미한다.

예시

```pdl
def {
    data_offset from 96 length 4
    payload from (*data_offset * 32) to end
}
```

만약 `data_offset` 필드의 비트값이 `0101`이라면, `*data_offset`은 정수 `5`로 평가된다.

---

### 정의

정의는 `def` 블록을 통해 이루어지며, 블록 내 각 한 줄마다 의미를 가진다.

각 줄은 필드 이름과 `from`, `to`, `length` 키워드 등을 통해 구성된다.

#### 길이 기반 정의

예시: 변수 `seq`를 32번 비트부터 시작하여 32개의 비트로 정의한다.

```pdl
def {
    seq from 32 length 32  # or `seq from 32;32`
}
```

이는 전체 비트 중 `[32] [33] [34] ... [62] [63]`을 의미한다.

#### 범위 기반 정의

예시: 변수 `src`를 0번 비트부터 시작하여 15번 비트까지로 정의한다.

```pdl
def {
    src from 0 to 15  # or `src from 0:15`
}
```

이는 전체 비트 중 `[0] [1] [2] ... [14] [15]`를 의미한다.

#### `to end`

`to end`는 시작 위치부터 입력 패킷의 마지막 비트까지를 의미한다.

예시: 변수 `payload`를 200번 비트부터 마지막까지로 정의

```pdl
def {
    payload from 200 to end
}
```

이는 다음과 동일한 의미를 가진다.

```text
payload = packet[200:]
```

#### `length remaining`

`remaining`은 현재 필드 시작 위치부터 패킷 끝까지 남아있는 모든 비트를 의미한다.

예시

```pdl
def {
    payload from 200 length remaining
}
```

이는 다음과 같은 의미이다.

```text
remaining = packet_total_bits - 200
```

따라서 아래 두 문장은 같은 의미로 취급할 수 있다.

```pdl
payload from 200 to end
payload from 200 length remaining
```

#### 다른 변수 기반 정의

다른 변수나 필드 값을 활용하여 필드의 길이를 정할 수 있다.

예시: `data_offset`에 할당되는 값을 정수로 해석하여, 그 값을 바탕으로 `options` 필드의 길이를 지정

```pdl
def {
    data_offset from 96  length 4
    options     from 160 length (*data_offset * 32 - 160)
}
```

필드 값 참조는 `*field_name` 형태를 사용한다.

```pdl
*data_offset
```

반면 `data_offset`만 사용하면 필드 이름 또는 필드 정의 자체를 가리키는 것으로 해석될 수 있으므로, 값 참조에는 반드시 `*`를 사용한다.

---

### 출력

출력은 `out` 블록을 통해 이루어지며, 블록 내 각 한 줄마다 의미를 가진다.

`out`은 [정의](#정의)에서 정의한 각 필드를 JSON 등과 어떻게 대응시킬지 정의한다.

각 줄은 변수 이름, JSON 이름, Format Type 키워드 등을 통해 구성된다.

예시: 변수 `seq`와 JSON `"sequence_number"` 필드를 대응시키고, 값은 10진수 숫자 포맷으로 해석 및 대입한다.

```pdl
out json {
    seq sequence_number DEC
}
```

출력 예시

```json
{
    "sequence_number": 11
}
```

JSON 이름 부분을 조정하여 JSON 구조를 설계할 수도 있다.

* `seq sequence.number DEC` => `{"sequence":{"number":11}}`
* `seq sequence[0].number DEC` => `{"sequence":[{"number":11}]}`
* `seq [0] DEC` => `[11]`

동일한 출력 대상에 여러 `out json` 블록이 존재할 경우, 출력 결과는 하나의 JSON 객체에 병합된다.

```pdl
out json {
    src source_port DEC
}

out json {
    dst destination_port DEC
}
```

출력

```json
{
    "source_port": 80,
    "destination_port": 443
}
```

---

### Endian

PDL은 Endian을 지원한다.

`set mode`는 값 해석 방식에 영향을 준다.

```pdl
set mode BIG_ENDIAN
```

또는

```pdl
set mode LITTLE_ENDIAN
```

각 블록은 자신의 위 가장 가까운 `set mode`의 Endian 방식을 따른다.

만약 `value`가 `0b11001010`이라면, `BIG_ENDIAN` 기준 비트 접근은 다음과 같다.

* `value<0>` -> 1
* `value<1>` -> 1
* `value<2>` -> 0
* `value<3>` -> 0
* ...
* `value<6>` -> 1
* `value<7>` -> 0

즉, `BIG_ENDIAN` 모드에서 `value<0>`은 최상위 비트(MSB)를 의미한다.

`LITTLE_ENDIAN` 모드에서는 구현 정책에 따라 `value<0>`을 최하위 비트(LSB)로 해석할 수 있다. MVP에서는 혼동을 줄이기 위해 다음 규칙을 권장한다.

```text
BIG_ENDIAN    => value<0> = MSB
LITTLE_ENDIAN => value<0> = LSB
```

---

### Format Type

| Format Type | Description                  | Example Bits                      | JSON Output                     |
| ----------- | ---------------------------- | --------------------------------- | ------------------------------- |
| `DEC`       | Unsigned Decimal Integer     | `00001110`                        | `14`                            |
| `HEX`       | Unsigned Hexadecimal Integer | `00001110`                        | `"0E"`                          |
| `BIN`       | Unsigned Binary Integer      | `00001110`                        | `"00001110"`                    |
| `BOOL`      | Boolean Type, only 1 bit     | `0`, `1`                          | `false`, `true`                 |
| `ASCII`     | ASCII Code Encoding          | `01100011 01100001`               | `"ca"`                          |
| `UTF8`      | UTF-8 Encoding String        | `01101000 01101001`               | `"hi"`                          |
| `UTF16`     | UTF-16 Encoding String       | `00000000 01101000 00000000 01101001` | `"hi"`                      |
| `BASE64`    | Base64 Representation        | `01100011 01100001`               | `"Y2E="`                        |
| `BASE58`    | Base58 Representation        | `01100011 01100001`               | `"4fzk"`                        |
| `IP4`       | IPv4 Type Bytes              | `11000000 10101000 00000000 00000001` | `"192.168.0.1"`            |
| `IP6`       | IPv6 Type Bytes              | `00000000...00000001` (128 bits)  | `"::1"`                         |
| `MAC`       | MAC Address Type Bytes       | `10101010 10111011 11001100 11011101 11101110 11111111` | `"AA:BB:CC:DD:EE:FF"` |

---

### Bracket Out

`out`은 값 매핑을 통해 출력 형태를 결정할 수도 있다.

#### 단일 비트 Flag 기반 출력

```pdl
out json {
    flags<6> syn {
        0 : false
        1 : true
    }
}
```

위 정의는 `flags` 필드의 6번 비트를 읽고, 그 값이 `0`이면 `false`, `1`이면 `true`로 출력한다.

출력 예시

```json
{
    "syn": true
}
```

#### 다중 비트 값 매핑

```pdl
out json {
    l3_protocol_type protocol {
        1  : "ICMP"
        6  : "TCP"
        17 : "UDP"
    }
}
```

출력 예시

```json
{
    "protocol": "TCP"
}
```

#### 매핑 실패

매핑 블록에 해당 값이 없을 경우 기본 동작은 원본 값을 출력하는 것이다.

예시

```pdl
out json {
    l3_protocol_type protocol {
        1  : "ICMP"
        6  : "TCP"
        17 : "UDP"
    }
}
```

만약 `l3_protocol_type = 132`라면 다음과 같이 출력한다.

```json
{
    "protocol": 132
}
```

---

## 오류 조건

PDL 파서는 다음 상황을 오류로 처리할 수 있다.

* 정의되지 않은 필드를 참조하는 경우
* `BOOL` 포맷에 2비트 이상 필드를 사용하는 경우
* `from` 또는 `length` 계산 결과가 음수인 경우
* 필드 범위가 입력 패킷 길이를 초과하는 경우
* 동일한 JSON 경로에 서로 다른 값이 중복 출력되는 경우
* `*field_name`을 사용했지만 해당 필드가 아직 파싱되지 않은 경우
* Bracket Out 매핑 키의 타입이 필드 값과 호환되지 않는 경우

---

## TCP 정의 예시

TCP 정의 예시는 별도 파일 `eg_tcp.pdl`에 작성한다.

```markdown
[TCP 정의 예시](./eg_tcp.pdl)
```

---

## 향후 확장 예정

### 조건부 필드

조건에 따라 서로 다른 필드 정의를 적용할 수 있다.

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

### 배열

반복 필드를 정의한다.

```pdl
array options count option_count {
    option_kind   length 8
    option_length length 8
}
```

### 구조체

재사용 가능한 하위 구조를 정의한다.

```pdl
struct Peer {
    id   length 256
    ip   length 32
    port length 16
}
```

### Include

다른 PDL 파일을 포함한다.

```pdl
include "./tcp.pdl"
include "./ipv4.pdl"
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
