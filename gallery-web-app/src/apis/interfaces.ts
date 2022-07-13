export interface IFileInfo {
  ID: number;
  Name: string;
  Thumbnail: string;
  ThumbnailHeight: number;
  ThumbnailWidth: number;
  RelativePath: string;
}
export interface IDirectoryInfo {
  ID: number;
  Name: string
  RelativePath: string;
  ParentDirectoryId: number;
}
export interface IMachine {
  ID: number;
  Name: string;
  IsOnline: boolean;
}
export interface IRootDirectoryInfo extends IDirectoryInfo{
  IsRootDirectory: boolean;
  Machine: IMachine
}