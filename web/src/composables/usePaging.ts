import type { Ref } from "vue";

export async function resetPageAndLoad(pageRef: Ref<number>, loader: () => Promise<void>) {
  pageRef.value = 1;
  await loader();
}
