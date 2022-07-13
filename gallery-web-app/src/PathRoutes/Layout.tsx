import * as React from "react";
import { NavLink, Outlet } from "react-router-dom";
import $style from "./Layout.module.less";
export interface ILayoutProps {}

export function Layout(props: ILayoutProps) {
  return (
    <div className={$style['layout-frame']}>
      <nav className={$style['layout-navbar']}>
        <NavLink to="/">MachineList</NavLink>
        {/* <NavLink to="/gallery">Gallery</NavLink> */}
        <NavLink to="/about"> About</NavLink>
      </nav>
      <Outlet />
    </div>
  );
}
