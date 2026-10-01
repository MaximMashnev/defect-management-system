import { Button, Container, IconButton, InputBase, styled } from "@mui/material";
import SearchIcon from '@mui/icons-material/Search';
import bgImage from "../../shared/assets/imgs/bgImage.jpg";
import { useState } from "react";
import { PATHS } from "../../shared/consts";
import { Link, useNavigate } from "react-router-dom";

const MainPageContainer = styled(Container)({
  position: "relative", 
  display: "flex", 
  alignItems: "center", 
  justifyContent: "center", 
  height: "100%",
  padding: 0,
  flexDirection: "column",
  gap: "1rem",
  '@media(min-width: 600px)' : {
    p: 0,
    maxWidth: "100%"
  },
  '&::before': {
    content: '""',
    position: "absolute",
    top: 0,
    left: 0,
    width: "100%",
    height: "100%",
    backgroundImage: `url(${bgImage})`,
    backgroundSize: "cover",
    backgroundPosition: "center",
    filter: "blur(2px)",
    zIndex: -1,
  }  
})

const MainApplicationForm = styled('form')({
  display: 'flex', 
  alignItems: 'center', 
  width: 400, 
  zIndex: 1,
  backgroundColor: "white",
  borderRadius: "4px"
})

export default function MainPage () {
  const [search, setSearch] = useState("");
  const navigate = useNavigate();

  const handleSearch = (e: React.FormEvent) => {
      e.preventDefault();
      if (search.trim()) {
        navigate(`/${PATHS.APPLICATION_TRACKING}/${search}`);
      }
  };

  return (
    <MainPageContainer>
      <MainApplicationForm onSubmit={handleSearch}>
        <InputBase
          sx={{ ml: 1, flex: 1, }}
          placeholder="Найти заявку"
          value={search}
          onChange={e => setSearch(e.target.value)}
        />
        <IconButton type="submit" sx={{ p: '10px' }} title="Найти заявку">
          <SearchIcon />
        </IconButton>
      </MainApplicationForm>
      <Button variant="contained" component={Link} to={PATHS.NEW_APPLICATION} fullWidth sx={{width: 400}}>
        Создать заявку
      </Button>
    </MainPageContainer>
  );
}