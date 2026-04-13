import { fileListApi } from "@/api/modules/disk";
import type { UserFile } from "@/types/api";
import { isFolder } from "@/utils/diskView";

const DEFAULT_CHUNK_SIZE = 200;

// The backend list endpoint is paginated; this helper hydrates full folder content for client-side operations.
export async function fetchAllFilesByParentId(parentId: number) {
  let currentPage = 1;
  let loaded: UserFile[] = [];
  let totalCount = 0;

  do {
    const response = await fileListApi({
      id: parentId,
      page: currentPage,
      size: DEFAULT_CHUNK_SIZE,
    });

    const list = response.list ?? [];
    loaded = loaded.concat(list);
    totalCount = response.count ?? 0;
    currentPage += 1;

    if (list.length === 0) {
      break;
    }
  } while (loaded.length < totalCount);

  return loaded;
}

export async function fetchFolderListByParentId(parentId: number) {
  const loaded = await fetchAllFilesByParentId(parentId);
  return loaded.filter((item) => isFolder(item));
}
