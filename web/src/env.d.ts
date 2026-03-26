/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_COS_BUCKET_URL?: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
