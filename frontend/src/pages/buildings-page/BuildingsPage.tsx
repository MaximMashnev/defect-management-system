import { 
  DataGrid, 
  getGridDateOperators, 
  GridActionsCell, 
  GridActionsCellItem, 
  GridColDef,
  GridFilterInputValue, 
  GridFilterModel,
  GridPaginationModel, 
  GridSortModel, 
  GridToolbarContainer, 
  GridToolbarFilterButton,
} from '@mui/x-data-grid';
import Paper from '@mui/material/Paper';
import buildingsStore from '../../entities/Buildings/model/store';
import { useEffect, useState } from 'react';
import { observer } from "mobx-react-lite"
import { Button, styled} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import BuildingFormDialog from '../../features/BuildingForm/ui/BuildingFormDialog';
import DeleteIcon from '@mui/icons-material/DeleteOutlined';
import EditIcon from '@mui/icons-material/Edit';
import RefreshIcon from '@mui/icons-material/Refresh';
import { BuildingEntity } from '../../entities/Buildings';
import { GridSortItem } from '@mui/x-data-grid/models/gridSortModel';

const customFilter = [
  {
    label: 'Равно',
    value: "",
    InputComponent: GridFilterInputValue,
    getApplyFilterFn: () => {return null},
  },
  {
    label: 'Содержит',
    value: '*',
    InputComponent: GridFilterInputValue,
    getApplyFilterFn: () => {return null},
  },
];

const PaperGrid = styled(Paper)({
  height: "auto",
  width: '100%' 
})

const BuildingsPage = observer(() => {

  const [openDialog, setOpenDialog] = useState(false);
  const [editData, setEditData] = useState<BuildingEntity | null>(null);
  const [paginationModel, setPaginationModel] = useState<GridPaginationModel>(buildingsStore.pagination);
  const [sortModel, setSortModel] = useState<GridSortItem | null>(buildingsStore.sort);
  const [filterModel, setFilterModel] = useState<GridFilterModel | null>(buildingsStore.filter);

  useEffect(() => {
    buildingsStore.setFilter(filterModel);
    buildingsStore.setPagination(paginationModel);
    buildingsStore.setSort(sortModel);
    buildingsStore.getBuildingsForTable();
  }, [paginationModel, sortModel, filterModel])

  const CustomToolbar = () => {
    return (
      <GridToolbarContainer>
        <GridToolbarFilterButton />
        <Button variant="text" startIcon={<AddIcon />} onClick={handleOpenDialog}>
          Новый объект
        </Button>
        <Button variant="text" startIcon={<RefreshIcon />} onClick={() => buildingsStore.getBuildingsForTable()}>
          Обновить
        </Button>
      </GridToolbarContainer>
    );
  }

  const columns: GridColDef<BuildingEntity>[] = [
    { 
      field: 'id', 
      headerName: 'ID', 
      width: 70,
      filterable: false
    },
    { 
      field: 'name', 
      headerName: 'Название', 
      type: 'longText',
      flex: 1,
      filterOperators: customFilter
    },
    { 
      field: 'address', 
      headerName: 'Адрес',
      type: 'longText',
      flex: 1,
      filterOperators: customFilter
    },
    {
      field: 'dateRegistration',
      headerName: 'Дата регистрации',
      flex: 1,
      type: "date",
      valueGetter: (value) => {return new Date(value)},
      filterOperators: getGridDateOperators().filter(
        (operator) => operator.value === 'is'
      ),
    },
    {
      field: 'numberApplications',
      headerName: 'Количество поданных заявок',
      type: 'number',
      flex: 1,
      filterable: false
    },
    {
      field: 'actions',
      headerName: 'Действия',
      type: "actions",
      renderCell: (params) => (
        <GridActionsCell {...params}>
          <GridActionsCellItem
            icon={<EditIcon />}
            label="Edit"
            onClick={() => handleEditClick(params.row)}
            color="inherit"
          />
          <GridActionsCellItem
            icon={<DeleteIcon />}
            label="Delete"
            onClick={() => handleDeleteClick(params.row)}
            color="inherit"
          />
        </GridActionsCell>
      ),
    }
  ];

  const handleOpenDialog = () => {
    setEditData(null);
    setOpenDialog(true);
  }

  const handleEditClick = (building: BuildingEntity) => {
    setEditData(building);
    setOpenDialog(true);
  }

  const handleDeleteClick = async (building: BuildingEntity) => {
    try {
      await buildingsStore.deleteBuilding(building);
      window.alert("Объект был удален.")
    }
    catch (error) {
      window.alert("Ошибка при удалении объекта.")
    }
  }  

  return (
    <>
      <PaperGrid>
        <DataGrid
            rows={buildingsStore.buildings}
            columns={columns}
            initialState={{ pagination: {paginationModel} }}
            pageSizeOptions={[5, 10]}
            paginationMode='server'
            onPaginationModelChange={(model: GridPaginationModel) => setPaginationModel(model)}
            rowCount={buildingsStore.meta?.total_items ?? 0}
            sortingMode="server"
            onSortModelChange={([model]: GridSortModel) => setSortModel(model)}
            filterMode='server'
            onFilterModelChange={(model: GridFilterModel) => setFilterModel(model)}
            sx={{ border: 0 }}
            showToolbar
            slots={{ toolbar: CustomToolbar }}
        />      
      </PaperGrid>
      <BuildingFormDialog open={openDialog} data={editData} onClose={() => setOpenDialog(false)}/>    
    </>
  )
})

export default BuildingsPage;