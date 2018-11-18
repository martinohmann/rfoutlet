import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';

import GithubLink from './GithubLink';
import Group from './Group';
import { apiRequest } from '../util';
import config from '../config';

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
    apiRequest('GET', '/status')
      .then(groups => this.setState({ groups }))
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
        <div className={classes.container}>
          {groups.map(group =>
            <Group key={group.id} {...group} />
          )}
        </div>
      </div>
    );
  }
}

Root.propTypes = {
 classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(Root);
