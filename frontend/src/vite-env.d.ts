/// <reference types="vite/client" />

interface ImportMetaEnv {
    readonly VITE_AUTH_SECRET: string
}

interface ImportMeta {
    readonly env: ImportMetaEnv
}
