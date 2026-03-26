import SparkMD5 from 'spark-md5';

export async function calcFileMd5(file: File, chunkSize = 2 * 1024 * 1024): Promise<string> {
  const spark = new SparkMD5.ArrayBuffer();
  const chunks = Math.ceil(file.size / chunkSize);

  for (let index = 0; index < chunks; index += 1) {
    const start = index * chunkSize;
    const end = Math.min(start + chunkSize, file.size);
    const buffer = await file.slice(start, end).arrayBuffer();
    spark.append(buffer);
  }

  return spark.end();
}
