import { create } from 'zustand';

interface GlobalStore {
  images: Record<string, string>;
  setImages: (images: Record<string, string>) => void;
}

export const useGlobalStore = create<GlobalStore>((set) => ({
  images: {},
  setImages: (images: Record<string, string>) => set({ images }),
}));