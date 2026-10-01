import { Typography } from '@mui/material';
import { PieChart, PieValueType } from '@mui/x-charts';
import { BarChart } from '@mui/x-charts/BarChart';
import { observer } from "mobx-react-lite";
import { useEffect, useMemo } from 'react';
import buildingsStore from '../../entities/Buildings/model/store';
import applicationsStore from '../../entities/Application/model/store';
import { ApplicationStatus } from '../../entities/Application/model/types';

interface BarChartData {
    label: string[];
    num: number[];
}

const DashboardPage = observer(() => {
    useEffect(() => {
        Promise.all([
            buildingsStore.getBuildings(),
            applicationsStore.getApplications()          
        ]);
    }, [])

    const statByStatus = useMemo((): PieValueType[] => {
        const appsStatusData: PieValueType[] = [];
        Object.values(ApplicationStatus).forEach(v => {
            const filterStatus = applicationsStore.applications.filter(a => a.status === v);
            appsStatusData.push({label: v, value: filterStatus.length});
        });
        return appsStatusData;
    }, [applicationsStore.applications]);

    const statByBuilding = useMemo((): BarChartData => {
        const buildingsData: BarChartData = {label: [], num: []};     
        buildingsStore.buildings.forEach(b => {
            buildingsData.label.push(b.name);
            buildingsData.num.push(b.numberApplications);
        })
        return buildingsData;
    }, [buildingsStore.buildings]);

    const statByTime = useMemo((): BarChartData => {
        const timeData: BarChartData = {label: [], num: []};
        applicationsStore.applications.forEach(app => {
            const dataHours = new Date(app.dateSubmission).getHours().toString();
            const indexTime = timeData.label.findIndex(l => l === dataHours);
            if (indexTime !== -1) {
                timeData.num[indexTime] += 1;
            }
            else {
                timeData.label.push(dataHours);
                timeData.num.push(1);
            }
        })
        return timeData;
    }, [applicationsStore.applications]);

    return (
        <>
            <Typography component="span" variant="h6">
                Количество заявок по статусам
            </Typography>
            <PieChart
                series={[
                    {
                        data: statByStatus,
                        highlightScope: { fade: 'global', highlight: 'item' },
                        faded: { innerRadius: 30, additionalRadius: -30, color: 'gray' },
                    },
                ]}
                height={200}
                width={200}
            />
            <Typography component="span" variant="h6">
                Статистика заявок по объектам
            </Typography>
            <BarChart
                xAxis={[{data: statByBuilding.label, height: 28}]}
                series={[{data: statByBuilding.num}]}
                height={300}
            />
            <Typography component="span" variant="h6">
                Тренды подачи заявок по времени
            </Typography>
            <BarChart
                xAxis={[{id: 'barCategories', data: statByTime.label, height: 28}]}
                series={[{data: statByTime.num}]}
                height={300}
            />
        </>
    )
})

export default DashboardPage;