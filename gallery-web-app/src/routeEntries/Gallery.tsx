import * as React from "react";
import PhotoAlbum, {  Photo } from "react-photo-album";
import {  useParams } from "react-router-dom";
import { getDirStructure,  IGetDirSturectureResponse } from "src/apis/getDirStructure";
import {  IDirectoryInfo, IFileInfo } from "src/apis/interfaces";
import DocumentPng from "src/assets/document.png";
import FolderPng from 'src/assets/folder.png';
import $style from './Gallery.module.less'
export const DEFAULT_PAGE_SIZE = 30;
export interface IGalleryPageProps {}
interface IFolderButtonProps {
  directory: IDirectoryInfo,
  setCurrentDirId: (value: number) => void
}
export function FolderButton(props: IFolderButtonProps) {
const{directory, setCurrentDirId} = props;
const onClick = React.useCallback(() => {
  setCurrentDirId(directory.ID);
},[directory.ID, setCurrentDirId])
return <button className={$style['subdir-btn']} onClick={onClick}>
  <img src={FolderPng} alt={`folder ${directory.ID}`}/>
  <p>[{directory.ID}]{directory.Name}</p>
</button>
}
const breakpoints = [2400, 1080, 640, 384, 256, 128, 96, 64, 48];

export function GalleryPage(props: IGalleryPageProps) {
  const { rootDirId } = useParams();
  const [pageIndex] = React.useState(0);
  const [photos, setPhotos] = React.useState<Photo[]>([]);
  const [subDirs, setSubDirs] = React.useState<IDirectoryInfo[]>([]);
  const [currentDirId, setCurrentDirId] = React.useState(Number(rootDirId ));
  const [dirData, setDirData] = React.useState<null | IGetDirSturectureResponse>(null);
  React.useEffect(() => {
    const controller = new AbortController();
    getDirStructure({
      DirId: currentDirId,
      PageOffset: pageIndex,
      PageSize:DEFAULT_PAGE_SIZE,
    }, {
      signal: controller.signal
    }).then((resp) => {
          setDirData(resp);
      })
      .catch((reason) => {
        if (reason.code === "ERR_CANCELED") {
          return;
        }
      });
    return () => controller.abort();
  }, [currentDirId, pageIndex]);
  React.useEffect(() => {
    if (dirData && Array.isArray(dirData.SubFiles)) {
      const newPhotos: Photo[] = dirData.SubFiles.map((f: IFileInfo) => {
        return {
          src: f.Thumbnail || DocumentPng,
          width: f.ThumbnailWidth,
          height: f.ThumbnailHeight,
          title: f.Name,
          key: String(f.ID),
          images: breakpoints.map((breakpoint) => {
            const height = Math.round((f.ThumbnailHeight / f.ThumbnailWidth) * breakpoint);
            return {
                src: f.Thumbnail || DocumentPng,
                width: breakpoint,
                height,
            };
          })
        };
      });
      console.log(newPhotos)
      setSubDirs(dirData.SubDirs);
      setPhotos(newPhotos);
    }
  }, [dirData]);
  const gotoParentDir = React.useCallback(() => {
    if (dirData?.Info.ParentDirectoryId)
    {setCurrentDirId(dirData.Info.ParentDirectoryId);}
  },[dirData?.Info.ParentDirectoryId]);
  const onPhotoClick = React.useCallback((_: any, photo: Photo) => {
    window.open(`/api/file/${photo.key}`)
  }, [])
  return (
    <div className={$style['gallery-frame']}>
      <div>
        <p>Current Path: <span>
          {dirData?.Info?.RelativePath}
          </span>
        </p>
        {dirData?.Info.ParentDirectoryId ? <button onClick={gotoParentDir} >
          Go to parent folder
            [{dirData?.Info.ParentDirectoryId}]...</button> : null}
      </div>
      <div className={$style['subdir-group']}>
        {subDirs.map(dir => <FolderButton key={`folder_${dir.ID}`} directory={dir} setCurrentDirId={setCurrentDirId}/>)}
      </div>
      <PhotoAlbum padding={4} layout="rows" onClick={onPhotoClick} photos={photos}/>
    </div>
  );
}
