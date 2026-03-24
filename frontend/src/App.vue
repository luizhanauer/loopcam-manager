<script setup lang="ts">
import { ref } from 'vue';
import { StartCamera, StopCamera, SelectVideoFile, SystemSetup } from '../wailsjs/go/main/App';

interface CameraState {
  devicePath: string;
  videoPath: string;
  isRunning: boolean;
  error: string;
}

const isSetupDone = ref(false);
const isSettingUp = ref(false);
const cameraCount = ref(2);
const setupError = ref('');

const cameras = ref<CameraState[]>([]);

const handleSetup = async () => {
  setupError.value = '';
  isSettingUp.value = true;
  
  try {
    const devices = await SystemSetup(cameraCount.value);
    devices.sort();
    
    cameras.value = devices.map(device => ({
      devicePath: device,
      videoPath: '',
      isRunning: false,
      error: ''
    }));
    
    isSetupDone.value = true;
  } catch (err: any) {
    setupError.value = "Erro ao configurar Kernel: " + err;
  } finally {
    isSettingUp.value = false;
  }
};

const handleSelectVideo = async (cam: CameraState) => {
  cam.error = '';
  try {
    const path = await SelectVideoFile();
    if (path) {
      cam.videoPath = path;
    }
  } catch (err: any) {
    cam.error = "Erro ao selecionar arquivo: " + err;
  }
};

const handleStart = async (cam: CameraState) => {
  cam.error = '';
  try {
    await StartCamera(cam.devicePath, cam.videoPath);
    cam.isRunning = true;
  } catch (err: any) {
    cam.error = err;
  }
};

const handleStop = async (cam: CameraState) => {
  try {
    await StopCamera(cam.devicePath);
    cam.isRunning = false;
  } catch (err: any) {
    cam.error = err;
  }
};
</script>

<template>
  <main class="min-h-screen bg-[#0b1121] text-slate-300 font-sans selection:bg-blue-500/30 relative overflow-hidden">
    <div class="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-blue-600/10 blur-[120px] rounded-full pointer-events-none"></div>
    <div class="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] bg-emerald-600/10 blur-[120px] rounded-full pointer-events-none"></div>

    <div class="max-w-5xl mx-auto px-6 py-10 relative z-10 space-y-10">
      
      <header class="flex flex-col md:flex-row md:items-end justify-between border-b border-slate-700/50 pb-6 gap-4">
        <div>
          <h1 class="text-4xl font-extrabold tracking-tight bg-gradient-to-r from-blue-400 to-emerald-400 bg-clip-text text-transparent">
            LoopCam Manager
          </h1>
          <p class="text-slate-400 mt-2 text-sm font-medium">Orquestrador de Webcams Virtuais V4L2 para Broadcast</p>
        </div>
        <div class="flex items-center space-x-2 text-xs font-mono bg-slate-800/50 border border-slate-700/50 px-3 py-1.5 rounded-lg text-slate-400">
          <span class="w-2 h-2 rounded-full bg-emerald-500/80"></span>
          <span>Kernel Sync Ready</span>
        </div>
      </header>

      <div v-if="isSetupDone" class="bg-blue-900/10 border border-blue-500/20 backdrop-blur-sm rounded-2xl p-5 shadow-lg">
        <h3 class="text-blue-300 font-semibold mb-3 flex items-center text-sm">
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
          Guia de Operação Segura
        </h3>
        <ul class="grid grid-cols-1 md:grid-cols-3 gap-4 text-xs text-slate-400">
          <li class="flex items-start">
            <span class="flex-shrink-0 w-5 h-5 flex items-center justify-center bg-blue-500/20 text-blue-300 rounded-full mr-2 font-bold">1</span>
            <span>Atribua um vídeo MP4 a uma câmera virtual.</span>
          </li>
          <li class="flex items-start">
            <span class="flex-shrink-0 w-5 h-5 flex items-center justify-center bg-blue-500/20 text-blue-300 rounded-full mr-2 font-bold">2</span>
            <span>Inicie a transmissão <b>antes</b> de abrir o OBS/Navegador.</span>
          </li>
          <li class="flex items-start">
            <span class="flex-shrink-0 w-5 h-5 flex items-center justify-center bg-blue-500/20 text-blue-300 rounded-full mr-2 font-bold">3</span>
            <span>Pare o loop aqui somente após fechar o uso da câmera no outro app.</span>
          </li>
        </ul>
      </div>

      <section v-if="!isSetupDone" class="bg-slate-800/40 backdrop-blur-md rounded-2xl p-8 shadow-2xl border border-slate-700/50 text-center max-w-lg mx-auto mt-16">
        <div class="w-16 h-16 bg-slate-700/50 rounded-2xl flex items-center justify-center mx-auto mb-6 shadow-inner">
          <svg class="w-8 h-8 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path></svg>
        </div>
        <h2 class="text-2xl font-semibold text-white mb-2">Alocação de Hardware</h2>
        <p class="text-sm text-slate-400 mb-8">Defina quantas webcams virtuais independentes o Kernel do Ubuntu deve instanciar.</p>
        
        <div class="flex items-center justify-center space-x-4 mb-8">
          <button @click="cameraCount > 1 && cameraCount--" class="w-10 h-10 rounded-lg bg-slate-700/50 hover:bg-slate-600 flex items-center justify-center text-xl transition-colors disabled:opacity-30" :disabled="isSettingUp">-</button>
          <input 
            v-model.number="cameraCount" 
            type="number" 
            min="1" max="10"
            class="w-20 bg-slate-900/50 border border-slate-600 rounded-xl px-2 py-3 text-center text-2xl font-bold text-white focus:ring-2 focus:ring-blue-500 focus:outline-none shadow-inner"
            :disabled="isSettingUp"
          />
          <button @click="cameraCount < 10 && cameraCount++" class="w-10 h-10 rounded-lg bg-slate-700/50 hover:bg-slate-600 flex items-center justify-center text-xl transition-colors disabled:opacity-30" :disabled="isSettingUp">+</button>
        </div>

        <button 
          @click="handleSetup"
          :disabled="isSettingUp"
          class="w-full bg-gradient-to-r from-blue-600 to-blue-500 hover:from-blue-500 hover:to-blue-400 text-white font-semibold py-3.5 px-6 rounded-xl transition-all shadow-lg hover:shadow-blue-500/25 focus:ring-2 focus:ring-offset-2 focus:ring-offset-slate-900 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center"
        >
          <span v-if="isSettingUp" class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin mr-3"></span>
          <span>{{ isSettingUp ? 'Autenticando e Injetando Módulo...' : 'Inicializar Kernel (pkexec)' }}</span>
        </button>

        <div v-if="setupError" class="mt-6 bg-red-500/10 border border-red-500/20 text-red-400 px-4 py-3 rounded-xl text-xs text-left shadow-inner flex items-start">
          <svg class="w-4 h-4 mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
          <span>{{ setupError }}</span>
        </div>
      </section>

      <div v-else class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        
        <section 
          v-for="cam in cameras" 
          :key="cam.devicePath"
          class="bg-slate-800/40 backdrop-blur-md rounded-2xl border transition-all duration-300 flex flex-col overflow-hidden"
          :class="cam.isRunning ? 'border-emerald-500/30 shadow-[0_0_15px_rgba(16,185,129,0.1)]' : 'border-slate-700/50 hover:border-slate-600/80 shadow-lg'"
        >
          <div class="bg-slate-800/80 px-5 py-4 border-b border-slate-700/50 flex justify-between items-center">
            <h2 class="text-base font-semibold text-white flex items-center space-x-3">
              <div class="relative flex h-3 w-3">
                <span v-if="cam.isRunning" class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                <span class="relative inline-flex rounded-full h-3 w-3" :class="cam.isRunning ? 'bg-emerald-500' : 'bg-slate-600'"></span>
              </div>
              <span>Webcam Virtual</span>
            </h2>
            <div class="flex items-center space-x-2">
              <span class="text-blue-400 font-mono text-xs bg-blue-900/20 px-2.5 py-1 rounded-md border border-blue-500/20 shadow-inner">
                {{ cam.devicePath }}
              </span>
            </div>
          </div>

          <div class="p-5 space-y-5 flex-1 flex flex-col justify-between">
            <div>
              <label class="block text-xs font-medium text-slate-400 mb-1.5 uppercase tracking-wider">Origem de Vídeo</label>
              <div class="flex gap-2">
                <div class="relative flex-1">
                  <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <svg class="w-4 h-4 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 4v16M17 4v16M3 8h4m10 0h4M3 12h18M3 16h4m10 0h4M4 20h16a1 1 0 001-1V5a1 1 0 00-1-1H4a1 1 0 00-1 1v14a1 1 0 001 1z"></path></svg>
                  </div>
                  <input 
                    v-model="cam.videoPath" 
                    type="text" 
                    placeholder="Selecione o arquivo .mp4..."
                    class="w-full bg-slate-900/50 border border-slate-600/80 rounded-lg pl-9 pr-3 py-2.5 text-sm text-slate-200 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-shadow disabled:opacity-50 disabled:cursor-not-allowed shadow-inner"
                    :disabled="cam.isRunning"
                  />
                </div>
                <button 
                  @click="handleSelectVideo(cam)"
                  :disabled="cam.isRunning"
                  title="Explorar arquivos"
                  class="bg-slate-700 hover:bg-slate-600 text-slate-200 px-4 py-2.5 rounded-lg text-sm transition-colors border border-slate-600 disabled:opacity-50 disabled:cursor-not-allowed focus:outline-none focus:ring-2 focus:ring-slate-500 font-medium flex items-center justify-center"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 19a2 2 0 01-2-2V7a2 2 0 012-2h4l2 2h4a2 2 0 012 2v1M5 19h14a2 2 0 002-2v-5a2 2 0 00-2-2H9a2 2 0 00-2 2v5a2 2 0 01-2 2z"></path></svg>
                </button>
              </div>
            </div>

            <div v-if="cam.error" class="bg-red-500/10 border border-red-500/20 text-red-400 px-3 py-2.5 rounded-lg text-xs break-words shadow-inner flex items-start">
               <svg class="w-4 h-4 mr-2 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
              {{ cam.error }}
            </div>

            <div class="pt-2">
              <button 
                v-if="!cam.isRunning"
                @click="handleStart(cam)"
                class="w-full bg-slate-700 hover:bg-blue-600 text-white font-semibold py-2.5 px-4 rounded-lg text-sm transition-all duration-200 border border-slate-600 hover:border-blue-500 focus:ring-2 focus:ring-blue-500 focus:outline-none flex items-center justify-center"
              >
                <svg class="w-4 h-4 mr-2" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM9.555 7.168A1 1 0 008 8v4a1 1 0 001.555.832l3-2a1 1 0 000-1.664l-3-2z" clip-rule="evenodd"></path></svg>
                Iniciar Transmissão
              </button>
              
              <button 
                v-else
                @click="handleStop(cam)"
                class="w-full bg-red-500/10 hover:bg-red-600 text-red-400 hover:text-white font-semibold py-2.5 px-4 rounded-lg text-sm transition-all duration-200 border border-red-500/30 hover:border-red-500 focus:ring-2 focus:ring-red-500 focus:outline-none flex items-center justify-center group"
              >
                <svg class="w-4 h-4 mr-2 text-red-400 group-hover:text-white transition-colors" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8 7a1 1 0 00-1 1v4a1 1 0 001 1h4a1 1 0 001-1V8a1 1 0 00-1-1H8z" clip-rule="evenodd"></path></svg>
                Parar Transmissão
              </button>
            </div>
          </div>
        </section>

      </div>
      
    </div>
  </main>
</template>