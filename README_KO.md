# Go Container System

[![License](https://img.shields.io/badge/license-BSD%203--Clause-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.25+-blue.svg)](https://golang.org/dl/)

> **언어:** [English](README.md) | **한국어**

## 개요

Go Container System은 Go를 위한 고성능 타입 안전 컨테이너 프레임워크로, 메시징 시스템 및 범용 애플리케이션을 위한 포괄적인 데이터 관리 기능을 제공합니다. 이것은 [C++ container_system](https://github.com/kcenon/container_system)의 Go 구현으로, Go의 장점을 활용하면서 동일한 기능을 제공하도록 설계되었습니다.

## 특징

### 핵심 기능

- **타입 안전 값 시스템**: 컴파일 타임 타입 체크를 제공하는 15가지 내장 값 타입
  - null, bool, short, ushort, int, uint, long, ulong, llong, ullong
  - float32, float64, bytes, string, container
- **메시지 컨테이너**: 헤더를 지원하는 완전한 기능의 메시지 컨테이너
  - 서브 ID를 포함한 소스/타겟 ID
  - 메시지 타입 및 버전 추적
  - 다중 직렬화 형식
- **직렬화**: 다중 형식 지원
  - 문자열 기반 직렬화
  - 바이트 배열 직렬화
  - JSON 변환
  - XML 변환
- **값 연산**: 풍부한 값 연산 세트
  - 타입 변환
  - 자식 값 관리
  - 이름으로 값 쿼리
- **컨테이너 연산**: 포괄적인 컨테이너 관리
  - 헤더 조작
  - 값 추가/제거/쿼리
  - 컨테이너 복사 (값 포함/제외)
  - 응답 메시지를 위한 헤더 스왑
- **플루언트 빌더 API**: 가독성 높은 컨테이너 생성을 위한 ContainerBuilder 패턴
  - source, target, type, values를 위한 체이닝 메서드
  - 선택적 스레드 안전 모드
- **의존성 주입 지원**: DI 프레임워크를 위한 표준 인터페이스 및 프로바이더
  - 손쉬운 모킹 및 테스트를 위한 ContainerFactory 인터페이스
  - 자동 와이어링을 위한 Google Wire 프로바이더 세트
  - Uber Dig 및 기타 DI 컨테이너와 호환

## 설치

```bash
go get github.com/kcenon/go_container_system
```

## 빠른 시작

### 간단한 값 생성

```go
import (
    "github.com/kcenon/go_container_system/container/core"
    "github.com/kcenon/go_container_system/container/values"
)

// 다양한 값 타입 생성
boolVal := values.NewBoolValue("enabled", true)
intVal := values.NewInt32Value("count", 42)
stringVal := values.NewStringValue("message", "안녕하세요!")

// 타입 변환
if val, err := intVal.ToInt32(); err == nil {
    fmt.Printf("값: %d\n", val)
}
```

### 메시지 컨테이너 생성

```go
// 전체 헤더를 포함한 컨테이너 생성
container := core.NewValueContainerFull(
    "client_app", "instance_1",  // 소스
    "server_api", "v2",           // 타겟
    "user_registration",          // 메시지 타입
)

// 값 추가
container.AddValue(values.NewStringValue("username", "alice"))
container.AddValue(values.NewInt32Value("age", 30))
container.AddValue(values.NewStringValue("email", "alice@example.com"))

// 직렬화
serialized, _ := container.Serialize()
jsonStr, _ := container.ToJSON()
xmlStr, _ := container.ToXML()
```

### 플루언트 빌더 API 사용

```go
import "github.com/kcenon/go_container_system/container/messaging"

// 플루언트 빌더 패턴으로 컨테이너 생성
container, err := messaging.NewContainerBuilder().
    WithSource("client", "1").
    WithTarget("server", "main").
    WithType("request").
    WithValues(
        values.NewStringValue("action", "login"),
        values.NewStringValue("user", "alice"),
    ).
    WithThreadSafe(true).
    Build()
```

### 의존성 주입 사용

```go
import "github.com/kcenon/go_container_system/container/di"

// ContainerFactory 직접 사용
factory := di.NewContainerFactory()
container := factory.NewContainer()
container = factory.NewContainerWithType("request")
builder := factory.NewBuilder()

// Google Wire와 함께 사용
// wire.go
//go:build wireinject
// +build wireinject

package main

import (
    "github.com/google/wire"
    "github.com/kcenon/go_container_system/container/di"
)

func InitializeApp() (*App, error) {
    wire.Build(di.ProviderSet, NewApp)
    return nil, nil
}

// Uber Dig와 함께 사용
container := dig.New()
container.Provide(di.NewContainerFactory)
```

### 컨테이너 값 작업

```go
// 중첩 구조 생성
userData := values.NewContainerValue("user",
    values.NewStringValue("name", "Bob"),
    values.NewInt32Value("age", 25),
)

// 부모 컨테이너에 추가
container.AddValue(userData)

// 값 검색
name := container.GetValue("name", 0)
if str, err := name.ToString(); err == nil {
    fmt.Printf("이름: %s\n", str)
}
```

## 아키텍처

### 패키지 구조

```
go_container_system/
├── container/
│   ├── core/           # 핵심 타입 및 인터페이스
│   │   ├── value_types.go   # 값 타입 열거형
│   │   ├── value.go         # Value 인터페이스 및 기본 구현
│   │   └── container.go     # ValueContainer 구현
│   ├── di/             # 의존성 주입 지원
│   │   ├── provider.go      # ContainerFactory 인터페이스 및 프로바이더
│   │   └── wire.go          # Google Wire 프로바이더 세트
│   ├── messaging/      # 플루언트 빌더 API
│   │   └── builder.go       # ContainerBuilder 구현
│   └── values/         # 구체적인 값 구현
│       ├── bool_value.go
│       ├── numeric_value.go
│       ├── string_value.go
│       ├── bytes_value.go
│       └── container_value.go
├── examples/           # 사용 예제
├── tests/             # 테스트 스위트
└── README.md
```

### 값 타입 계층 구조

```
Value (인터페이스)
├── BaseValue (기본 구현)
│   ├── BoolValue
│   ├── Int16Value, UInt16Value
│   ├── Int32Value, UInt32Value
│   ├── Int64Value, UInt64Value
│   ├── Float32Value, Float64Value
│   ├── StringValue
│   ├── BytesValue
│   └── ContainerValue
└── ValueContainer (메시지 컨테이너)
```

## 값 타입

### 숫자 타입

| 타입 | Go 타입 | 크기 | 설명 |
|------|---------|------|-------------|
| ShortValue | int16 | 2 바이트 | 16비트 부호 있는 정수 |
| UShortValue | uint16 | 2 바이트 | 16비트 부호 없는 정수 |
| IntValue | int32 | 4 바이트 | 32비트 부호 있는 정수 |
| UIntValue | uint32 | 4 바이트 | 32비트 부호 없는 정수 |
| LongValue | int32 | 4 바이트 | 32비트 부호 있는 정수 (호환성) |
| ULongValue | uint32 | 4 바이트 | 32비트 부호 없는 정수 (호환성) |
| LLongValue | int64 | 8 바이트 | 64비트 부호 있는 정수 |
| ULLongValue | uint64 | 8 바이트 | 64비트 부호 없는 정수 |
| FloatValue | float32 | 4 바이트 | 32비트 부동소수점 |
| DoubleValue | float64 | 8 바이트 | 64비트 부동소수점 |

### 기타 타입

- **BoolValue**: 불리언 (true/false)
- **StringValue**: UTF-8 문자열
- **BytesValue**: 바이너리 데이터
- **ContainerValue**: 자식 값을 가진 중첩 컨테이너
- **NullValue**: 빈/null 값

## 사용 사례

- **메시지 전달**: IPC를 위한 구조화된 메시지 컨테이너
- **네트워크 프로토콜**: 네트워크 통신을 위한 바이너리 직렬화
- **설정**: 유연한 설정 데이터 구조
- **데이터 교환**: 언어 간 데이터 직렬화
- **API 통신**: REST API를 위한 JSON/XML 직렬화

## C++ 버전과의 호환성

이 Go 구현은 C++ container_system과 동일한 기능을 제공합니다:

### 동일한 기능
- ✅ 동일한 의미를 가진 15가지 값 타입
- ✅ 헤더를 지원하는 값 컨테이너
- ✅ 문자열 및 바이트 배열 직렬화
- ✅ XML 및 JSON 변환
- ✅ 컨테이너 복사 연산
- ✅ 헤더 스왑 기능
- ✅ 이름 및 인덱스로 값 쿼리

### Go 특화 개선 사항
- 🔹 더 나은 타입 안전성을 위한 인터페이스 기반 설계
- 🔹 Go 관용구를 사용한 에러 처리 (에러 반환)
- 🔹 가비지 컬렉션 (수동 메모리 관리 불필요)
- 🔹 Go 관례를 사용한 단순화된 API

### 아직 구현되지 않음
- ⏳ MessagePack 직렬화 (계획됨)
- ⏳ 파일 로드/저장 연산 (계획됨)
- ⏳ 뮤텍스를 사용한 스레드 안전 연산 (계획됨)
- ⏳ 메모리 풀 최적화 (Go에서는 불필요)

## 라이선스

이 프로젝트는 BSD 3-Clause 라이선스로 라이선스가 부여됩니다. 자세한 내용은 LICENSE 파일을 참조하세요.

## 기여

기여를 환영합니다! Pull Request를 자유롭게 제출해 주세요.

## 작성자

**kcenon**
- 이메일: kcenon@naver.com
- GitHub: [@kcenon](https://github.com/kcenon)

## 감사의 글

- C++ [container_system](https://github.com/kcenon/container_system) 기반
- 메시징 시스템 생태계와의 호환성을 위해 설계됨
