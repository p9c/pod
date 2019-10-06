package css

func HELPERS() string {
	return `
.lineBottom{
box-shadow: inset 0 -1px 0 rgba(0, 0, 0, .13);
}
.lineRight{
box-shadow: inset -1px 0 0 rgba(0, 0, 0, .12);
}
html, body{
width:100%;
height:100vh;
margin:0;
padding:0;
}
#x{
position:fixed;
margin:0;
padding:0;
top:0;
left:0;
right:0;
bottom:0;
overflow:hidden;
background-color:var(--dark);
}
#display{
position:relative;
display:block;
width:100%;
height:100%;
}

.rwrap{
width: 100%!important;
height:100%!important;
}
.cwrap{
position: relative;
width: 100%!important;
height:100%!important;
}


.posRel{
position: relative!important;
}
.posAbs{
position: absolute!important;
}

.flx{
display:flex!important;
}
.flc{
flex-direction: column;
}
.fii{
flex:1;
}


.marZ{
margin:0!important;
}
.padZ{
padding:0!important;
}


.lsn{
list-style:none!important;
}

.overHid{
overflow:hidden!important;
}

.marginTop{
margin-top:30px;
}

.justifyCenter{
justify-content: center;
}
.justifyBetween{
justify-content: space-between;
}

.justifyEvenly{
justify-content: space-evenly;
}

.itemsCenter{
align-items: center;
}

.noMargin{
margin:0!important;
}

.noPadding{
padding:0!important;
}


.baseMargin{
margin:15!important;
}






.hide{
display:none;
}

.Logo{
justify-content:center;
align-items:center;
background:var(--dark);
}
.Logo svg{
width:auto;
height:48px;
}



.logofill{
fill:var(--light);
background:var(--light);
}
`
}
