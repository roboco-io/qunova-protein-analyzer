# PRD: 단백질 구조 분석 시스템 (Protein Structure Analyzer)

## 1. 제품 개요

### 1.1 제품명
단백질 구조 분석 시스템 (Protein Structure Analyzer)

### 1.2 제품 목적
과학자와 연구원들이 단백질 구조 파일(PDB, CIF, mmCIF)을 업로드하여 구조를 자동으로 분석하고, 체인, 도메인, 리간드별로 데이터를 분리하여 활용할 수 있도록 지원하는 웹 기반 분석 도구

### 1.3 핵심 가치
- **자동화된 구조 분석**: 복잡한 단백질 구조를 자동으로 파싱하고 분석
- **다층 데이터 분리**: 체인, 도메인, 리간드별 세분화된 데이터 추출
- **직관적인 인터페이스**: 드래그 앤 드롭 방식의 간편한 파일 업로드
- **즉시 활용 가능한 결과**: JSON 형식의 구조화된 데이터 내보내기

## 2. 주요 기능

### 2.1 파일 업로드
- **지원 형식**: PDB, CIF, mmCIF 파일
- **업로드 방식**: 
  - 드래그 앤 드롭
  - 파일 선택 버튼
  - 데모 모드 (샘플 데이터로 시연)

### 2.2 구조 분석
- **자동 파싱**: 업로드된 파일의 구조 자동 분석
- **실시간 진행률**: 분석 진행 상황 시각적 표시
- **분석 요소**:
  - 총 원자 수 계산
  - 체인 식별 및 분류
  - 도메인 매핑
  - 리간드 검출

### 2.3 결과 시각화
- **요약 대시보드**: 주요 지표 한눈에 확인
  - 총 원자 수
  - 체인 수
  - 도메인 수
  - 리간드 수

### 2.4 세부 분석 결과

#### 체인별 분석
- 체인 ID 및 타입 표시 (Protein/DNA/Water)
- 각 체인의 잔기(residue) 수
- 각 체인의 원자 수
- 색상 코딩으로 타입 구분

#### 도메인별 분석
- 도메인 명칭
- 잔기 범위
- 연관 체인
- 포함 원자 수

#### 리간드 분석
- 리간드 명칭
- 결합 부위 정보
- 원자 수

### 2.5 데이터 내보내기
- **형식**: JSON
- **내보내기 옵션**:
  - 체인 데이터
  - 도메인 데이터
  - 리간드 데이터
  - 전체 데이터

### 2.6 연계 기능
- 분자 역학 시뮬레이션으로 데이터 전송 기능

## 3. 사용자 경험

### 3.1 워크플로우
1. **업로드 단계**: 파일 선택/드래그 또는 데모 시작
2. **처리 단계**: 자동 분석 진행 (진행률 표시)
3. **결과 단계**: 분석 결과 확인 및 데이터 활용

### 3.2 UI/UX 특징
- **진행 단계 인디케이터**: 현재 진행 상황 시각적 표시
- **반응형 디자인**: 모바일/태블릿/데스크톱 대응
- **색상 테마**: 청색 계열 그라데이션 배경
- **아이콘 활용**: Lucide 아이콘으로 직관적 정보 전달

## 4. 기술 사양

### 4.1 프론트엔드
- **프레임워크**: React
- **스타일링**: Tailwind CSS
- **아이콘**: Lucide React
- **상태 관리**: React Hooks (useState, useRef)

### 4.2 주요 컴포넌트
- 파일 업로드 핸들러
- 진행률 트래커
- 데이터 시각화 카드
- 내보내기 기능

## 5. 목표 사용자

### 5.1 주요 사용자
- 생명과학 연구원
- 제약회사 연구개발팀
- 대학 연구실
- 바이오인포매틱스 전문가

### 5.2 사용 시나리오
- 단백질 구조 연구
- 약물 표적 분석
- 구조 기반 약물 설계
- 교육 및 시연 목적

## 6. 성공 지표

- **사용성**: 5분 이내 첫 분석 완료
- **정확도**: 구조 파싱 정확도 99% 이상
- **성능**: 10MB 파일 30초 이내 처리
- **만족도**: 사용자 만족도 4.5/5.0 이상

## 7. 향후 개발 계획

### 단기 (3개월)
- 실제 파일 파싱 엔진 통합
- 3D 구조 시각화 추가
- 배치 처리 기능

### 중기 (6개월)
- AI 기반 구조 예측 통합
- 협업 기능 추가
- API 제공

### 장기 (12개월)
- 클라우드 기반 대용량 처리
- 머신러닝 기반 인사이트 제공
- 다른 분석 도구와의 통합

---

# 단백질 구조 분석 시스템 - 기술 문서

## 📋 목차
1. [개요](#개요)
2. [설치 및 설정](#설치-및-설정)
3. [컴포넌트 구조](#컴포넌트-구조)
4. [주요 기능 구현](#주요-기능-구현)
5. [데이터 구조](#데이터-구조)
6. [API 참조](#api-참조)
7. [사용 예제](#사용-예제)
8. [문제 해결](#문제-해결)

## 개요

단백질 구조 분석 시스템은 React 기반의 웹 애플리케이션으로, 단백질 구조 파일을 업로드하여 분석하고 구조화된 데이터를 추출할 수 있는 도구입니다.

### 주요 특징
- ✅ PDB, CIF, mmCIF 파일 지원
- ✅ 드래그 앤 드롭 업로드
- ✅ 실시간 분석 진행률 표시
- ✅ 체인, 도메인, 리간드별 데이터 분리
- ✅ JSON 형식 데이터 내보내기

## 설치 및 설정

### 필수 요구사항
- Node.js 14.0 이상
- npm 또는 yarn
- React 17.0 이상

### 설치
```bash
# 의존성 설치
npm install react lucide-react
npm install -D tailwindcss

# 컴포넌트 추가
cp protein-structure-analyzer.tsx src/components/
```

### 사용
```jsx
import ProteinStructureAnalyzer from './components/protein-structure-analyzer';

function App() {
  return <ProteinStructureAnalyzer />;
}
```

## 컴포넌트 구조

### 주요 상태 관리
```jsx
const [file, setFile] = useState(null);           // 업로드된 파일
const [analysisStep, setAnalysisStep] = useState('upload'); // 현재 단계
const [progress, setProgress] = useState(0);      // 진행률
const [structureData, setStructureData] = useState(null); // 분석 결과
```

### 컴포넌트 계층 구조
```
ProteinStructureAnalyzer
├── Header (제목 및 설명)
├── ProgressIndicator (진행 단계 표시)
└── MainContent
    ├── UploadSection (analysisStep === 'upload')
    ├── ProcessingSection (analysisStep === 'processing')
    └── ResultsSection (analysisStep === 'results')
        ├── SummaryCard
        ├── ChainsAnalysis
        ├── DomainsAnalysis
        ├── LigandsAnalysis
        └── NextSteps
```

## 주요 기능 구현

### 1. 파일 업로드 처리

```jsx
const handleFileUpload = (event) => {
  const uploadedFile = event.target?.files?.[0];
  const fileName = uploadedFile ? uploadedFile.name : 'demo_protein_1abc.pdb';
  setFile({ name: fileName });
  setAnalysisStep('processing');
  
  // 진행률 시뮬레이션
  let currentProgress = 0;
  const interval = setInterval(() => {
    currentProgress += Math.random() * 20;
    setProgress(Math.min(currentProgress, 100));
    
    if (currentProgress >= 100) {
      clearInterval(interval);
      setTimeout(() => {
        setStructureData({ ...mockStructureData, fileName });
        setAnalysisStep('results');
      }, 500);
    }
  }, 300);
};
```

### 2. 드래그 앤 드롭 구현

```jsx
const handleDragOver = (e) => {
  e.preventDefault();
};

const handleDrop = (e) => {
  e.preventDefault();
  const droppedFile = e.dataTransfer?.files?.[0];
  // 파일 처리 로직
};
```

### 3. 데이터 내보내기

```jsx
const exportData = (type, data) => {
  const blob = new Blob([JSON.stringify(data, null, 2)], { 
    type: 'application/json' 
  });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `${type}_${structureData.fileName.split('.')[0]}.json`;
  a.click();
};
```

## 데이터 구조

### 구조 데이터 스키마

```typescript
interface StructureData {
  fileName: string;
  totalAtoms: number;
  chains: Chain[];
  domains: Domain[];
  ligands: Ligand[];
}

interface Chain {
  id: string;
  residues: number;
  atoms: number;
  type: 'Protein' | 'DNA' | 'RNA' | 'Water';
}

interface Domain {
  name: string;
  range: string;
  chains: string[];
  atoms: number;
}

interface Ligand {
  name: string;
  atoms: number;
  binding_site: string;
}
```

### 예제 데이터

```json
{
  "fileName": "protein_1abc.pdb",
  "totalAtoms": 2847,
  "chains": [
    {
      "id": "A",
      "residues": 142,
      "atoms": 1123,
      "type": "Protein"
    }
  ],
  "domains": [
    {
      "name": "N-terminal domain",
      "range": "1-67",
      "chains": ["A"],
      "atoms": 534
    }
  ],
  "ligands": [
    {
      "name": "ATP",
      "atoms": 31,
      "binding_site": "Chain A:45-52"
    }
  ]
}
```

## API 참조

### 주요 함수

#### `handleFileUpload(event)`
파일 입력 이벤트를 처리하고 분석을 시작합니다.

**매개변수:**
- `event`: 파일 입력 이벤트 객체

#### `startDemo()`
데모 데이터로 분석 프로세스를 시뮬레이션합니다.

#### `exportData(type, data)`
지정된 데이터를 JSON 파일로 내보냅니다.

**매개변수:**
- `type`: 데이터 유형 ('chains' | 'domains' | 'ligands')
- `data`: 내보낼 데이터 객체

#### `resetAnalysis()`
모든 상태를 초기화하고 새 분석을 준비합니다.

## 사용 예제

### 기본 사용

```jsx
import ProteinStructureAnalyzer from './ProteinStructureAnalyzer';

function App() {
  return (
    <div className="App">
      <ProteinStructureAnalyzer />
    </div>
  );
}
```

### 커스텀 설정

```jsx
// 실제 백엔드 API와 연동하는 경우
const CustomAnalyzer = () => {
  const analyzeProtein = async (file) => {
    const formData = new FormData();
    formData.append('file', file);
    
    const response = await fetch('/api/analyze', {
      method: 'POST',
      body: formData
    });
    
    return response.json();
  };
  
  // 컴포넌트 로직...
};
```

## 문제 해결

### 자주 발생하는 문제

#### 1. 파일 업로드가 작동하지 않음
- 파일 형식이 .pdb, .cif, .mmcif인지 확인
- 파일 크기가 브라우저 제한을 초과하지 않는지 확인

#### 2. 진행률이 멈춤
- 네트워크 연결 확인
- 브라우저 콘솔에서 오류 메시지 확인

#### 3. 데이터 내보내기 실패
- 브라우저 다운로드 권한 확인
- 팝업 차단기 설정 확인

### 성능 최적화

```jsx
// 대용량 파일 처리를 위한 청크 처리
const processLargeFile = async (file) => {
  const chunkSize = 1024 * 1024; // 1MB
  const chunks = Math.ceil(file.size / chunkSize);
  
  for (let i = 0; i < chunks; i++) {
    const start = i * chunkSize;
    const end = Math.min(start + chunkSize, file.size);
    const chunk = file.slice(start, end);
    // 청크 처리...
  }
};
```

## 향후 개선사항

- [ ] 실제 PDB 파일 파싱 엔진 통합
- [ ] 3D 구조 시각화 (Three.js/Mol*)
- [ ] 서버 사이드 처리 구현
- [ ] 배치 파일 처리
- [ ] 결과 데이터 캐싱
- [ ] 고급 필터링 옵션

## 라이선스

MIT License

## 기여 가이드

Pull Request를 환영합니다. 주요 변경사항의 경우 먼저 이슈를 열어 논의해 주세요.