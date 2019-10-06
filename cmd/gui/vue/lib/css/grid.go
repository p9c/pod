package css

func GRID() string {
	return `
.grid-container {
    height: 100%;
    margin: 0;
}

.grid-container * {
    box-shadow: inset 0 0 0 1px rgba(0, 0, 0, .3);
    position: relative;
}

.grid-container *:after {
    position: absolute;
    top: 0;
    left: 0;
}

.grid-container {
    display: grid;
    grid-template-columns: 60px 1fr;
    grid-template-rows: 60px 1fr;
    grid-template-areas: "Logo Header" "Sidebar Main";
}

.Logo {
    grid-area: Logo;
}

.Header {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr 60px 180px 120px 120px 60px;
    grid-template-rows: 1fr;
    grid-template-areas: "h1 h2 h3 h4 h5 h6 h7 h8";
    grid-gap: 5px;
    grid-area: Header;
}

.h1 {
    grid-area: h1;
}

.h2 {
    grid-area: h2;
}

.h3 {
    grid-area: h3;
}

.h4 {
    grid-area: h4;
}

.h5 {
    grid-area: h5;
}

.h6 {
    grid-area: h6;
}

.h7 {
    grid-area: h7;
}

.h8 {
    grid-area: h8;
}

.Sidebar {
    display: grid;
    grid-template-columns: 1fr;
    grid-template-rows: 600px 1fr 60px;
    grid-template-areas: "Nav" "Side" "Open";
    grid-area: Sidebar;
}

.Open {
    grid-area: Open;
}

.Nav {
    grid-area: Nav;
}

.Side {
    grid-area: Side;
}

.Main {
    padding: 15px;
    display: grid;
    grid-gap: 15px;
    grid-area: Main;
   grid-template-columns: 1fr 1fr 1fr 1fr 1fr 1fr 1fr 1fr 1fr;
  grid-template-rows: 1fr 1fr 1fr 1fr;
  grid-template-areas: "Balance Balance Balance Balance Send Send Send Send Send" "Txs Txs Txs Txs Txs Log Log NetHash NetHash" "Txs Txs Txs Txs Txs Log Log LocalHash LocalHash" "Info Info Time Time Time Log Log Status Status";
}

.Balance { grid-area: Balance; }

.Send { grid-area: Send; }

.NetHash { grid-area: NetHash; }

.LocalHash { grid-area: LocalHash; }

.Status { grid-area: Status; }

.Txs { grid-area: Txs; }

.Log { grid-area: Log; }

.Info { grid-area: Info; }

.Time { grid-area: Time; }

`
}
