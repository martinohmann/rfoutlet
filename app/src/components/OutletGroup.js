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
import { apiRequest, outletEnabled } from '../util'

const styles = {
  outletGroup: {
    paddingTop: 1,
    paddingBottom: 1,
    paddingRight: 6,
    background: grey[100],
  },
  buttonOn: {
    color: green[500],
  },
  buttonOff: {
    color: red[500],
  },
  outletGroupIdentifier: {
    flexGrow: 1,
    fontWeight: 700,
    color: grey[800],
  },
};

class OutletGroup extends React.Component {
  constructor(props, context) {
    super(props, context)

    this.registerOutlet = this.registerOutlet.bind(this);

    this.state = {
      outlets: {},
    }
  }

  registerOutlet(id, outlet) {
    const outlets = this.state.outlets;

    outlets[id] = outlet

    this.setState({ outlets: outlets });
  }

  updateOutletStates(outlets) {
    outlets.map(outlet => {
      if (undefined !== this.state.outlets[outlet.id]) {
        this.state.outlets[outlet.id].setState({
          isEnabled: outletEnabled(outlet),
        });
      }

      return outlet;
    });
  }

  handleButtonClick(action) {
    const data = {
      action: action,
      id: this.props.group.id
    };

    apiRequest('POST', '/outlet_group', data)
      .then(result => this.updateOutletStates(result.outlets))
      .catch(err => console.error(err));
  }

  render() {
    const { name, outlets } = this.props.group;
    const { classes } = this.props;

    return (
      <div>
        <List component="nav">
          <ListItem className={classes.outletGroup}>
            <ListItemText primary={name} disableTypography={true} className={classes.outletGroupIdentifier} />
            <div>
              <IconButton className={classes.buttonOff} onClick={(e) => this.handleButtonClick('off') }>
                <Icon>flash_off</Icon>
              </IconButton>
              <IconButton className={classes.buttonOn} onClick={(e) => this.handleButtonClick('on') }>
                <Icon>flash_on</Icon>
              </IconButton>
              <IconButton onClick={(e) => this.handleButtonClick('toggle') }>
                <Icon>swap_horiz</Icon>
              </IconButton>
            </div>
          </ListItem>
          {outlets.map(outlet =>
            <Outlet
              key={outlet.id}
              outlet={outlet}
              registerOutlet={this.registerOutlet}
            />
          )}
        </List>
      </div>
    );
  }
}

OutletGroup.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(OutletGroup);
