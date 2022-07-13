import * as React from "react";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import { GalleryPage } from "src/routeEntries/Gallery";
import { MachineList } from "src/routeEntries/MachineList";
import { Layout } from "./Layout";
export interface IPathRoutesProps {}

export function PathRoutes(props: IPathRoutesProps) {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Layout />} >
          <Route path="gallery/:rootDirId" element={<GalleryPage />} />
          <Route path="/" element={<MachineList />}/>
        </Route>
      </Routes>
    </BrowserRouter>
  );
}
