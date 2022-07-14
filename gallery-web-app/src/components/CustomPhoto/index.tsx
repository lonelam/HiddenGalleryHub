import path from 'path';
import * as React from 'react';
import { RenderPhoto } from 'react-photo-album';
import DocumentPng from "src/assets/document.png";
import { ICON_MAP } from './constants';

export const CustomPhoto : RenderPhoto = ({ layout, layoutOptions, imageProps: { alt, style, title, src, srcSet,...restImageProps } }) =>{
    const thumbnailSrc = React.useMemo(() => {
        if (src) {
            return src;
        }
        const lowerExt = path.extname(title || '').toLowerCase();
    if (ICON_MAP[lowerExt]) {
        return ICON_MAP[lowerExt].src
    }
    return DocumentPng;
    },[src, title])
  return (
    <div  style={{
        border: "2px solid #eee",
        borderRadius: "4px",
        boxSizing: "content-box",
        alignItems: "center",
        width: style?.width,
        padding: `${layoutOptions.padding - 2}px`,
        paddingBottom: 0,
    }}>
        <img alt={title} 
        style={{ ...style, width: src? "100%" : "50%", 
        padding: 0, margin: "0 auto" }} 
        src={thumbnailSrc}
        {...restImageProps}
        />
        <div
            style={{
                paddingTop: "8px",
                paddingBottom: "8px",
                overflow: "visible",
                whiteSpace: "nowrap",
                textAlign: "center",
            }}
        >
            {title}
        </div>
    </div>
  );
}
