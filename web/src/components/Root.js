import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';

import GithubLink from './GithubLink';
import GroupList from './GroupList';
import config from '../config';
import WebSocket from '../websocket';

const styles = theme => ({
  root: {
    flexGrow: 1,
  },
  title: {
    flexGrow: 1,
    color: theme.palette.common.white,
  },
  container: {
    marginTop: 64,
  },
});

class Root extends React.Component {
  state = {
    groups: [],
  }

  componentDidMount() {
    this.connectWebSocket();
    this.dispatchMessage({ type: 'status' });
  }

  connectWebSocket() {
    this.ws = new WebSocket(config.ws.url);
    this.ws.attachDefaultListeners();
    this.ws.onMessage(msg => {
      console.log("[ws recv]", msg);
      this.setState({ groups: msg });
    });
  }

  dispatchMessage = (msg) => {
    console.log("[ws send]", msg);
    this.ws.sendMessage(msg)
      .catch(err => console.error(err));
  }

  render() {
    const { classes } = this.props;
    const { groups } = this.state;

    return (
      <div className={classes.root}>
        <AppBar position="fixed">
          <Toolbar>
            <Typography variant="h6" className={classes.title}>
              {config.project.name}
            </Typography>
            <GithubLink url={config.project.url} />
          </Toolbar>
        </AppBar>
        <GroupList groups={groups} dispatchMessage={this.dispatchMessage} />
      </div>
    );
  }
}

Root.propTypes = {
 classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(Root);
