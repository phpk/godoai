function _writeStringToArray(aString: string, targetArray: Uint8Array, offset: number): void {
  for (let i = 0; i < aString.length; ++i)
    targetArray[offset + i] = aString.charCodeAt(i);
}

function _writeInt16ToArray(aNumber: number, targetArray: Uint8Array, offset: number): void {
  aNumber = Math.floor(aNumber);
  targetArray[offset + 0] = aNumber & 255;          // byte 1
  targetArray[offset + 1] = (aNumber >> 8) & 255;   // byte 2
}

function _writeInt32ToArray(aNumber: number, targetArray: Uint8Array, offset: number): void {
  aNumber = Math.floor(aNumber);
  targetArray[offset + 0] = aNumber & 255;          // byte 1
  targetArray[offset + 1] = (aNumber >> 8) & 255;   // byte 2
  targetArray[offset + 2] = (aNumber >> 16) & 255;  // byte 3
  targetArray[offset + 3] = (aNumber >> 24) & 255;  // byte 4
}

function _floatBits(f: number): number {
  const buf = new ArrayBuffer(4);
  (new Float32Array(buf))[0] = f;
  const bits = (new Uint32Array(buf))[0];
  // Return as a signed integer.
  return bits | 0;
}

function _writeAudioBufferToArray(
  audioBuffer: AudioBuffer,
  targetArray: Uint8Array,
  offset: number,
  bitDepth: number
): void {
  let index = 0, channel = 0;
  const length = audioBuffer.length;
  const channels = audioBuffer.numberOfChannels;
  let channelData, sample;

  for (index = 0; index < length; ++index) {
    for (channel = 0; channel < channels; ++channel) {
      channelData = audioBuffer.getChannelData(channel);

      if (bitDepth === 16) {
        sample = channelData[index] * 32768.0;
        sample = Math.min(Math.max(sample, -32768), 32767);
        _writeInt16ToArray(sample, targetArray, offset);
        offset += 2;
      } else if (bitDepth === 32) {
        sample = _floatBits(channelData[index]);
        _writeInt32ToArray(sample, targetArray, offset);
        offset += 4;
      } else {
        console.log('Invalid bit depth for PCM encoding.');
        return;
      }
    }
  }
}
async function _getAudioBuffer(blobData: any, contextOptions?: AudioContextOptions): Promise<AudioBuffer | null> {
  try {
    const audioContext = new AudioContext(contextOptions);
    audioContext.state === 'suspended' && audioContext.resume(); // Resume the context if suspended
    const arrayBuffer = await fetch(URL.createObjectURL(blobData)).then((res) => res.arrayBuffer());
    return await audioContext.decodeAudioData(arrayBuffer);
  } catch (error) {
    console.error("Error decoding audio data:", error);
    return null;
  }
}

export async function getWaveBlob(
  blobData: Blob | Blob[],
  as32BitFloat: boolean,
  contextOptions?: AudioContextOptions
): Promise<Blob> {
  const audioBuffer = await _getAudioBuffer(blobData, contextOptions);
  // 新增错误处理
  if (!audioBuffer) {
    console.error("Failed to retrieve audio buffer.");
    throw new Error("Audio buffer is null.");
  }

  const frameLength = audioBuffer.length;
  const numberOfChannels = audioBuffer.numberOfChannels;
  const sampleRate = audioBuffer.sampleRate;
  const bitsPerSample = as32BitFloat ? 32 : 16;
  const bytesPerSample = bitsPerSample / 8;
  const byteRate = sampleRate * numberOfChannels * bitsPerSample / 8;
  const blockAlign = numberOfChannels * bitsPerSample / 8;
  const wavDataByteLength = frameLength * numberOfChannels * bytesPerSample;
  const headerByteLength = 44;
  const totalLength = headerByteLength + wavDataByteLength;
  const waveFileData = new Uint8Array(totalLength);
  const subChunk1Size = 16;
  const subChunk2Size = wavDataByteLength;
  const chunkSize = 4 + (8 + subChunk1Size) + (8 + subChunk2Size);

  _writeStringToArray('RIFF', waveFileData, 0);
  _writeInt32ToArray(chunkSize, waveFileData, 4);
  _writeStringToArray('WAVE', waveFileData, 8);
  _writeStringToArray('fmt ', waveFileData, 12);

  _writeInt32ToArray(subChunk1Size, waveFileData, 16);
  _writeInt16ToArray(as32BitFloat ? 3 : 1, waveFileData, 20);
  _writeInt16ToArray(numberOfChannels, waveFileData, 22);
  _writeInt32ToArray(sampleRate, waveFileData, 24);
  _writeInt32ToArray(byteRate, waveFileData, 28);
  _writeInt16ToArray(blockAlign, waveFileData, 32);
  _writeInt32ToArray(bitsPerSample, waveFileData, 34);
  _writeStringToArray('data', waveFileData, 36);
  _writeInt32ToArray(subChunk2Size, waveFileData, 40);

  _writeAudioBufferToArray(audioBuffer, waveFileData, 44, bitsPerSample);

  return new Blob([waveFileData], { type: 'audio/wave' });
}
