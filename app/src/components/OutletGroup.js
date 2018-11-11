import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import IconButton from '@material-ui/core/IconButton';
import Icon from '@material-ui/core/Icon';
import red from '@material-ui/core/colors/red';
import green from '@material-ui/core/colors/green';
import grey from '@material-ui/core/colors/grey';

import Outlet from './Outlet';
import { apiRequest, outletEnabled } from '../util';

const styles = {
  group: {
    paddingTop: 1,
    paddingBottom: 1,
    paddingRight: 6,
    background: grey[100],
  },
  groupIdentifier: {
    flexGrow: 1,
    fontWeight: 700,
    color: grey[800],
  },
  buttonGroupOn: {
    color: green[500],
  },
  buttonGroupOff: {
    color: red[500],
  },
};

class OutletGroup extends React.Component {
  state = {
    outlets: [],
  }

  registerOutlet = (outlet) => {
    this.setState(state => {
      return { outlets: [...state.outlets, outlet] }
    });
  }

  updateOutletStates = outlets => {
    outlets.map((outlet, id) => {
      const o = this.state.outlets[id];

      if (undefined !== o) {
        o.setState(state => {
          return {
            enabled: outletEnabled(outlet),
            timeSwitch: { ...state.timeSwitch, enabled: false },
          };
        });
      }

      return outlet;
    });
  }

  handleButtonClick = action => event => {
    const { groupId } = this.props;

    apiRequest('POST', '/outlet_group', { groupId, action })
      .then(result => this.updateOutletStates(result.outlets))
      .catch(err => console.error(err));
  }

  render() {
    const { identifier, outlets } = this.props.attributes;
    const { classes, groupId } = this.props;

    return (
      <List component="nav">
        <ListItem className={classes.group}>
          <ListItemText primary={identifier} disableTypography={true} className={classes.groupIdentifier} />
          <div>
            <IconButton className={classes.buttonGroupOff} onClick={this.handleButtonClick('off')}>
              <Icon>flash_off</Icon>
            </IconButton>
            <IconButton className={classes.buttonGroupOn} onClick={this.handleButtonClick('on')}>
              <Icon>flash_on</Icon>
            </IconButton>
            <IconButton onClick={this.handleButtonClick('toggle')}>
              <Icon>swap_horiz</Icon>
            </IconButton>
          </div>
        </ListItem>
        {outlets.map((outlet, outletId) =>
          <Outlet
            key={outletId}
            outletId={outletId}
            groupId={groupId}
            attributes={outlet}
            registerOutlet={this.registerOutlet}
          />
        )}
      </List>
    );
  }
}

OutletGroup.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(OutletGroup);
