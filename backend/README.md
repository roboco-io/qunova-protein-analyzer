# Protein Structure Analyzer Backend

Go 기반의 단백질 구조 분석 백엔드 서비스입니다.

## 기능

- 🔬 PDB, CIF, mmCIF 파일 업로드 및 파싱
- 🧬 체인, 도메인, 리간드 분석
- 📊 구조 데이터 JSON 내보내기
- 🚀 RESTful API 제공
- 🐳 Docker 컨테이너 지원

## 기술 스택

- **언어**: Go 1.21+
- **웹 프레임워크**: Gin
- **파싱**: Python BioPython
- **컨테이너**: Docker & Docker Compose

## 빠른 시작

### 사전 요구사항

- Go 1.21 이상
- Python 3.9 이상
- Docker & Docker Compose
- Make (선택사항)

### 로컬 개발

1. 의존성 설치:
```bash
cd backend
make setup
```

2. 환경 변수 설정:
```bash
cp .env.example .env
```

3. 애플리케이션 실행:
```bash
make run
# 또는
go run .
```

### Docker로 실행

1. 전체 스택 실행:
```bash
docker-compose up
```

2. 백엔드만 실행:
```bash
cd backend
make docker-build
make docker-run
```

## API 엔드포인트

### Health Check
```
GET /health
```

### 파일 업로드
```
POST /api/v1/upload
Content-Type: multipart/form-data
Body: file (PDB/CIF/mmCIF)

Response:
{
  "fileId": "uuid",
  "fileName": "protein.pdb",
  "fileSize": 12345,
  "fileType": "pdb",
  "uploadedAt": "2025-01-18T..."
}
```

### 구조 분석 시작
```
POST /api/v1/analyze/{fileId}
Content-Type: application/json

Body (optional):
{
  "options": {
    "includeSequence": true,
    "calculateDomains": true,
    "identifyLigands": true,
    "resolution": "high"
  }
}

Response:
{
  "analysisId": "uuid",
  "status": "processing",
  "message": "Analysis started successfully"
}
```

### 분석 결과 조회
```
GET /api/v1/analysis/{analysisId}

Response:
{
  "id": "uuid",
  "status": "completed",
  "progress": 100,
  "structure": {
    "fileName": "protein.pdb",
    "totalAtoms": 2847,
    "chains": [...],
    "domains": [...],
    "ligands": [...]
  }
}
```

### 체인 데이터 조회
```
GET /api/v1/analysis/{analysisId}/chains
```

### 도메인 데이터 조회
```
GET /api/v1/analysis/{analysisId}/domains
```

### 리간드 데이터 조회
```
GET /api/v1/analysis/{analysisId}/ligands
```

### 데이터 내보내기
```
GET /api/v1/export/{analysisId}
GET /api/v1/export/{analysisId}/{dataType}

dataType: chains | domains | ligands
```

## 프로젝트 구조

```
backend/
├── main.go              # 애플리케이션 진입점
├── go.mod               # Go 모듈 정의
├── internal/
│   ├── handlers/        # HTTP 핸들러
│   ├── services/        # 비즈니스 로직
│   ├── models/          # 데이터 모델
│   └── middleware/      # 미들웨어
├── scripts/
│   └── parse_protein.py # Python 파서
├── uploads/             # 업로드된 파일 저장
├── Dockerfile           # Docker 이미지 정의
├── Makefile            # 빌드 자동화
└── .env.example        # 환경 변수 예제
```

## 개발

### 테스트 실행
```bash
make test
```

### 코드 포맷팅
```bash
make fmt
```

### 린팅
```bash
make lint
```

### Hot Reload 개발
```bash
make dev
```

## 환경 변수

| 변수 | 설명 | 기본값 |
|------|------|--------|
| PORT | 서버 포트 | 8080 |
| GIN_MODE | Gin 모드 (debug/release) | release |
| UPLOAD_MAX_SIZE | 최대 업로드 크기 | 104857600 |
| USE_MOCK_PARSER | 모의 파서 사용 여부 | true |

## 문제 해결

### Python 파서 오류
BioPython이 설치되지 않은 경우:
```bash
pip3 install biopython
```

### 포트 충돌
8080 포트가 사용 중인 경우 `.env` 파일에서 PORT 변경

### Docker 빌드 실패
Docker 데몬이 실행 중인지 확인:
```bash
docker ps
```

## 라이선스

MIT License