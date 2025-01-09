/// <reference types="vite/client" />
interface ImportMetaEnv {
    readonly VITE_TMP_ENV: string;
}

interface ImportMeta {
    readonly env: ImportMetaEnv;
}
