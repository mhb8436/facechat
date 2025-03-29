# FaceChat IM Server

FaceChat은 Go 언어로 작성된 고성능 인스턴트 메시징 서버 시스템입니다. 마이크로서비스 아키텍처를 기반으로 하며, WebSocket과 gRPC를 활용한 실시간 통신을 지원합니다.

## 주요 기능

- 실시간 메시지 전송 및 수신
- WebSocket 기반 양방향 통신
- gRPC를 통한 서비스 간 통신
- Redis를 활용한 메시지 큐잉
- SQLite/PostgreSQL 데이터베이스 지원
- etcd를 통한 서비스 디스커버리
- 분산 시스템 아키텍처

## 시스템 구성

프로젝트는 다음과 같은 주요 모듈들로 구성되어 있습니다:

- **api**: HTTP API 서버
- **connect**: WebSocket 연결 관리
- **message**: 메시지 처리 및 라우팅
- **sender**: 메시지 전송 서비스
- **db**: 데이터베이스 관리
- **proto**: Protocol Buffers 정의
- **tools**: 유틸리티 함수들

## 기술 스택

- Go 1.19
- Gin (Web Framework)
- gRPC
- WebSocket
- Redis
- SQLite/PostgreSQL
- etcd
- Protocol Buffers

## 시작하기

### 사전 요구사항

- Go 1.19 이상
- Redis
- etcd (선택사항)

### 설치

```bash
# 프로젝트 클론
git clone [repository-url]

# 의존성 설치
go mod download
```

### 실행

프로젝트는 4개의 주요 서비스로 구성되어 있으며, `r.sh` 스크립트를 통해 모든 서비스를 한 번에 실행할 수 있습니다:

```bash
# 실행 권한 부여
chmod +x r.sh

# 서비스 실행
./r.sh
```

또는 개별 서비스를 실행할 수 있습니다:

```bash
# API 서버
go run main.go -module api

# 메시지 서버
go run main.go -module message

# 연결 서버
go run main.go -module connect

# 전송 서버
go run main.go -module sender
```

## 프로젝트 구조

```
.
├── api/            # HTTP API 서버
├── config/         # 설정 파일
├── connect/        # WebSocket 연결 관리
├── db/             # 데이터베이스 관련
├── message/        # 메시지 처리
├── proto/          # Protocol Buffers 정의
├── public/         # 정적 파일
├── sender/         # 메시지 전송
└── tools/          # 유틸리티 함수
```

## 라이선스

이 프로젝트는 MIT 라이선스를 따릅니다. 자세한 내용은 [LICENSE](LICENSE) 파일을 참조하세요.

```
MIT License

Copyright (c) 2024 FaceChat

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## 기여하기

프로젝트에 기여하고 싶으시다면 Pull Request를 보내주세요.
