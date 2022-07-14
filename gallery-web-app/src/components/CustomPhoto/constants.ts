import Cr2Png from 'src/assets/cr2.png';
import PsdPng from 'src/assets/psd.png';
export const ICON_MAP = {
    '.cr2': {
        src: Cr2Png
    },
    '.psd': {
        src: PsdPng
    }
} as {[k: string]: {
    src: string;
}};