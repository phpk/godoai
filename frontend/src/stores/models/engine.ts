export const modelEngines = [
  {
    name: "ollama",
    cpp:"llama.cpp",
    needQuant : true
  },
  // {
  //   name: "embedding",
  //   cpp: "llama.cpp",
  //   needQuant: true
  // },
  {
    name: "sd",
    cpp: "stable-diffusion.cpp",
    needQuant: false
  },
  // {
  //   name: "rwkv",
  //   cpp: "rwkv.cpp",
  //   needQuant: true
  // },
  {
    name: "voice",
    cpp: "sherpa.cpp",
    needQuant: false
  }

]
