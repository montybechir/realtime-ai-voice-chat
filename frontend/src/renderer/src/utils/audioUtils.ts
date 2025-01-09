export const convertInt16ToFloat32 = (input: Int16Array): Float32Array => {
    const output = new Float32Array(input.length);
    for (let i = 0; i < input.length; i++) {
        const s = Math.max(-32768, Math.min(32767, input[i]));
        output[i] = s < 0 ? s / 32768 : s / 32767;
    }
    return output;
};
