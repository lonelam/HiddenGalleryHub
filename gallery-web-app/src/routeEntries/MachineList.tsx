import * as React from 'react';
import { getMachineList } from 'src/apis/getMachineList';
import {  IRootDirectoryInfo } from 'src/apis/interfaces';
import $style from './MachineList.module.less';
import NeuralPng from 'src/assets/neural.png'
import { useNavigate } from 'react-router-dom';
export interface IMachineListProps {
}
interface IMachineButton {
    machineDir: IRootDirectoryInfo
}
export function GoToMachineRootButton (props: IMachineButton) {
    const navigate = useNavigate()
    const {machineDir} = props; const onClick = React.useCallback(() => {
        navigate(`/gallery/${machineDir.ID}`);
    },[machineDir.ID, navigate]);
    return <button className={machineDir.Machine.IsOnline ? $style['offline-machine']: ''} onClick={onClick}>
        <img alt="the machine you connected" src={NeuralPng}/>
        <p>ID: {machineDir.Machine.ID}</p>
        <p>Name: {machineDir.Machine.Name}</p>
        <p>RootDir: {machineDir.Name}</p>
        </button>
}

export function MachineList (props: IMachineListProps) {
    const [machines, setMachines] = React.useState<IRootDirectoryInfo[]>([])
    React.useEffect(() => {
        getMachineList().then((resp) => {
            setMachines(resp.machines);
        })
    },[])
  return (
    <div className={$style['machine-list-wrapper']}>
      {
        machines.map(machineDir => 
           <GoToMachineRootButton key={machineDir.ID} machineDir={machineDir}/>
        )
      }
    </div>
  );
}
