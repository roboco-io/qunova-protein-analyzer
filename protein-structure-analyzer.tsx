import React, { useState, useRef } from 'react';
import { Upload, FileText, Atom, Layers, Download, Play, CheckCircle, AlertCircle, Zap } from 'lucide-react';

const ProteinStructureAnalyzer = () => {
  const [file, setFile] = useState(null);
  const [analysisStep, setAnalysisStep] = useState('upload'); // upload, processing, results
  const [progress, setProgress] = useState(0);
  const [structureData, setStructureData] = useState(null);
  const fileInputRef = useRef(null);

  // 모의 단백질 구조 데이터
  const mockStructureData = {
    fileName: 'protein_1abc.pdb',
    totalAtoms: 2847,
    chains: [
      { id: 'A', residues: 142, atoms: 1123, type: 'Protein' },
      { id: 'B', residues: 89, atoms: 712, type: 'Protein' },
      { id: 'C', residues: 45, atoms: 356, type: 'DNA' },
      { id: 'W', residues: 127, atoms: 381, type: 'Water' }
    ],
    domains: [
      { name: 'N-terminal domain', range: '1-67', chains: ['A'], atoms: 534 },
      { name: 'C-terminal domain', range: '68-142', chains: ['A'], atoms: 589 },
      { name: 'Binding domain', range: '1-89', chains: ['B'], atoms: 712 }
    ],
    ligands: [
      { name: 'ATP', atoms: 31, binding_site: 'Chain A:45-52' },
      { name: 'MG', atoms: 1, binding_site: 'Chain A:48' }
    ]
  };

  const handleFileUpload = (event) => {
    const uploadedFile = event.target?.files?.[0];
    
    // 실제 파일이 있으면 그 이름을 사용하고, 없으면 모의 파일 이름 사용
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

  const startDemo = () => {
    // 모의 파일 설정
    setFile({ name: 'demo_protein_1abc.pdb' });
    setAnalysisStep('processing');
    
    // 진행률 시뮬레이션
    let currentProgress = 0;
    const interval = setInterval(() => {
      currentProgress += Math.random() * 20;
      setProgress(Math.min(currentProgress, 100));
      
      if (currentProgress >= 100) {
        clearInterval(interval);
        setTimeout(() => {
          setStructureData({ ...mockStructureData, fileName: 'demo_protein_1abc.pdb' });
          setAnalysisStep('results');
        }, 500);
      }
    }, 300);
  };

  const handleDragOver = (e) => {
    e.preventDefault();
  };

  const handleDrop = (e) => {
    e.preventDefault();
    const droppedFile = e.dataTransfer?.files?.[0];
    
    // 실제 파일이 있으면 그 이름을 사용하고, 없으면 모의 파일 이름 사용
    const fileName = droppedFile ? droppedFile.name : 'dropped_protein_structure.pdb';
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

  const exportData = (type, data) => {
    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${type}_${structureData.fileName.split('.')[0]}.json`;
    a.click();
  };

  const resetAnalysis = () => {
    setFile(null);
    setAnalysisStep('upload');
    setProgress(0);
    setStructureData(null);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 p-6">
      <div className="max-w-6xl mx-auto">
        {/* 헤더 */}
        <div className="text-center mb-8">
          <div className="flex items-center justify-center mb-4">
            <Atom className="w-12 h-12 text-blue-600 mr-3" />
            <h1 className="text-4xl font-bold text-gray-800">단백질 구조 분석 시스템</h1>
          </div>
          <p className="text-xl text-gray-600">구조 파일을 업로드하여 분석하고 분리된 데이터를 생성하세요</p>
        </div>

        {/* 진행 단계 표시 */}
        <div className="flex items-center justify-center mb-8">
          <div className="flex items-center space-x-4">
            <div className={`flex items-center ${analysisStep === 'upload' ? 'text-blue-600' : 'text-green-600'}`}>
              <div className={`w-8 h-8 rounded-full flex items-center justify-center ${
                analysisStep === 'upload' ? 'bg-blue-100' : 'bg-green-100'
              }`}>
                {analysisStep === 'upload' ? '1' : <CheckCircle className="w-5 h-5" />}
              </div>
              <span className="ml-2 font-medium">파일 업로드</span>
            </div>
            <div className="w-12 h-0.5 bg-gray-300"></div>
            <div className={`flex items-center ${
              analysisStep === 'processing' ? 'text-blue-600' : 
              analysisStep === 'results' ? 'text-green-600' : 'text-gray-400'
            }`}>
              <div className={`w-8 h-8 rounded-full flex items-center justify-center ${
                analysisStep === 'processing' ? 'bg-blue-100' : 
                analysisStep === 'results' ? 'bg-green-100' : 'bg-gray-100'
              }`}>
                {analysisStep === 'results' ? <CheckCircle className="w-5 h-5" /> : '2'}
              </div>
              <span className="ml-2 font-medium">구조 분석</span>
            </div>
            <div className="w-12 h-0.5 bg-gray-300"></div>
            <div className={`flex items-center ${analysisStep === 'results' ? 'text-blue-600' : 'text-gray-400'}`}>
              <div className={`w-8 h-8 rounded-full flex items-center justify-center ${
                analysisStep === 'results' ? 'bg-blue-100' : 'bg-gray-100'
              }`}>
                3
              </div>
              <span className="ml-2 font-medium">결과 확인</span>
            </div>
          </div>
        </div>

        {/* 메인 컨텐츠 */}
        {analysisStep === 'upload' && (
          <div className="bg-white rounded-lg shadow-lg p-8">
            <div 
              className="border-2 border-dashed border-blue-300 rounded-lg p-12 text-center hover:border-blue-400 transition-colors cursor-pointer"
              onDragOver={handleDragOver}
              onDrop={handleDrop}
              onClick={startDemo}
            >
              <Upload className="w-16 h-16 text-blue-400 mx-auto mb-4" />
              <h3 className="text-xl font-semibold text-gray-700 mb-2">단백질 구조 파일 업로드</h3>
              <p className="text-gray-500 mb-4">PDB, CIF, mmCIF 파일을 드래그하여 놓거나 클릭하여 데모를 시작하세요</p>
              <button 
                onClick={startDemo}
                className="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition-colors mr-4"
              >
                데모 시작
              </button>
              <button 
                onClick={() => fileInputRef.current?.click()}
                className="bg-gray-600 text-white px-6 py-2 rounded-lg hover:bg-gray-700 transition-colors"
              >
                실제 파일 선택
              </button>
              <input
                ref={fileInputRef}
                type="file"
                accept=".pdb,.cif,.mmcif"
                onChange={handleFileUpload}
                className="hidden"
              />
            </div>
            
            <div className="mt-8 grid grid-cols-1 md:grid-cols-3 gap-6">
              <div className="text-center">
                <FileText className="w-8 h-8 text-green-600 mx-auto mb-2" />
                <h4 className="font-semibold text-gray-700">지원 형식</h4>
                <p className="text-sm text-gray-500">PDB, CIF, mmCIF</p>
              </div>
              <div className="text-center">
                <Zap className="w-8 h-8 text-orange-600 mx-auto mb-2" />
                <h4 className="font-semibold text-gray-700">빠른 분석</h4>
                <p className="text-sm text-gray-500">자동 구조 파싱</p>
              </div>
              <div className="text-center">
                <Layers className="w-8 h-8 text-purple-600 mx-auto mb-2" />
                <h4 className="font-semibold text-gray-700">다층 분리</h4>
                <p className="text-sm text-gray-500">체인, 도메인, 잔기별</p>
              </div>
            </div>
          </div>
        )}

        {analysisStep === 'processing' && (
          <div className="bg-white rounded-lg shadow-lg p-8">
            <div className="text-center">
              <div className="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-6">
                <Atom className="w-8 h-8 text-blue-600 animate-spin" />
              </div>
              <h3 className="text-2xl font-semibold text-gray-700 mb-4">구조 분석 중...</h3>
              <p className="text-gray-500 mb-6">단백질 구조를 파싱하고 분리하고 있습니다</p>
              
              <div className="w-full bg-gray-200 rounded-full h-4 mb-4">
                <div 
                  className="bg-blue-600 h-4 rounded-full transition-all duration-300"
                  style={{ width: `${progress}%` }}
                ></div>
              </div>
              <p className="text-sm text-gray-500">{Math.round(progress)}% 완료</p>
            </div>
          </div>
        )}

        {analysisStep === 'results' && structureData && (
          <div className="space-y-6">
            {/* 결과 요약 */}
            <div className="bg-white rounded-lg shadow-lg p-6">
              <div className="flex items-center justify-between mb-6">
                <div>
                  <h3 className="text-2xl font-semibold text-gray-800 mb-2">분석 완료</h3>
                  <p className="text-gray-600">파일: {file?.name || structureData.fileName}</p>
                </div>
                <button 
                  onClick={resetAnalysis}
                  className="bg-gray-100 text-gray-700 px-4 py-2 rounded-lg hover:bg-gray-200 transition-colors"
                >
                  새 분석 시작
                </button>
              </div>
              
              <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
                <div className="bg-blue-50 p-4 rounded-lg text-center">
                  <div className="text-2xl font-bold text-blue-600">{structureData.totalAtoms.toLocaleString()}</div>
                  <div className="text-sm text-gray-600">총 원자 수</div>
                </div>
                <div className="bg-green-50 p-4 rounded-lg text-center">
                  <div className="text-2xl font-bold text-green-600">{structureData.chains.length}</div>
                  <div className="text-sm text-gray-600">체인 수</div>
                </div>
                <div className="bg-purple-50 p-4 rounded-lg text-center">
                  <div className="text-2xl font-bold text-purple-600">{structureData.domains.length}</div>
                  <div className="text-sm text-gray-600">도메인 수</div>
                </div>
                <div className="bg-orange-50 p-4 rounded-lg text-center">
                  <div className="text-2xl font-bold text-orange-600">{structureData.ligands.length}</div>
                  <div className="text-sm text-gray-600">리간드 수</div>
                </div>
              </div>
            </div>

            {/* 체인 분석 */}
            <div className="bg-white rounded-lg shadow-lg p-6">
              <div className="flex items-center justify-between mb-4">
                <h4 className="text-xl font-semibold text-gray-800">체인별 분석</h4>
                <button 
                  onClick={() => exportData('chains', structureData.chains)}
                  className="flex items-center bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors"
                >
                  <Download className="w-4 h-4 mr-2" />
                  체인 데이터 내보내기
                </button>
              </div>
              <div className="space-y-3">
                {structureData.chains.map((chain, index) => (
                  <div key={index} className="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                    <div className="flex items-center">
                      <div className={`w-10 h-10 rounded-full flex items-center justify-center text-white font-bold mr-4 ${
                        chain.type === 'Protein' ? 'bg-blue-500' : 
                        chain.type === 'DNA' ? 'bg-green-500' : 'bg-gray-500'
                      }`}>
                        {chain.id}
                      </div>
                      <div>
                        <div className="font-semibold text-gray-800">체인 {chain.id} ({chain.type})</div>
                        <div className="text-sm text-gray-600">{chain.residues}개 잔기, {chain.atoms}개 원자</div>
                      </div>
                    </div>
                    <button className="text-blue-600 hover:text-blue-800 font-medium">
                      상세 보기
                    </button>
                  </div>
                ))}
              </div>
            </div>

            {/* 도메인 분석 */}
            <div className="bg-white rounded-lg shadow-lg p-6">
              <div className="flex items-center justify-between mb-4">
                <h4 className="text-xl font-semibold text-gray-800">도메인별 분석</h4>
                <button 
                  onClick={() => exportData('domains', structureData.domains)}
                  className="flex items-center bg-purple-600 text-white px-4 py-2 rounded-lg hover:bg-purple-700 transition-colors"
                >
                  <Download className="w-4 h-4 mr-2" />
                  도메인 데이터 내보내기
                </button>
              </div>
              <div className="space-y-3">
                {structureData.domains.map((domain, index) => (
                  <div key={index} className="p-4 bg-purple-50 rounded-lg">
                    <div className="flex items-center justify-between">
                      <div>
                        <div className="font-semibold text-gray-800">{domain.name}</div>
                        <div className="text-sm text-gray-600">
                          범위: {domain.range} | 체인: {domain.chains.join(', ')} | {domain.atoms}개 원자
                        </div>
                      </div>
                      <button className="text-purple-600 hover:text-purple-800 font-medium">
                        분리하기
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* 리간드 분석 */}
            <div className="bg-white rounded-lg shadow-lg p-6">
              <div className="flex items-center justify-between mb-4">
                <h4 className="text-xl font-semibold text-gray-800">리간드 분석</h4>
                <button 
                  onClick={() => exportData('ligands', structureData.ligands)}
                  className="flex items-center bg-orange-600 text-white px-4 py-2 rounded-lg hover:bg-orange-700 transition-colors"
                >
                  <Download className="w-4 h-4 mr-2" />
                  리간드 데이터 내보내기
                </button>
              </div>
              <div className="space-y-3">
                {structureData.ligands.map((ligand, index) => (
                  <div key={index} className="p-4 bg-orange-50 rounded-lg">
                    <div className="flex items-center justify-between">
                      <div>
                        <div className="font-semibold text-gray-800">{ligand.name}</div>
                        <div className="text-sm text-gray-600">
                          결합 부위: {ligand.binding_site} | {ligand.atoms}개 원자
                        </div>
                      </div>
                      <button className="text-orange-600 hover:text-orange-800 font-medium">
                        추출하기
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* 다음 단계 */}
            <div className="bg-gradient-to-r from-blue-600 to-indigo-600 rounded-lg shadow-lg p-6 text-white">
              <h4 className="text-xl font-semibold mb-4">다음 단계로 진행</h4>
              <p className="mb-6">분리된 구조 데이터를 다운로드하거나 다음 분석 단계로 전송하세요.</p>
              <div className="flex space-x-4">
                <button className="bg-white text-blue-600 px-6 py-2 rounded-lg hover:bg-gray-100 transition-colors font-medium">
                  전체 데이터 다운로드
                </button>
                <button className="bg-blue-500 text-white px-6 py-2 rounded-lg hover:bg-blue-400 transition-colors font-medium flex items-center">
                  <Play className="w-4 h-4 mr-2" />
                  분자 역학 시뮬레이션으로 전송
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default ProteinStructureAnalyzer;