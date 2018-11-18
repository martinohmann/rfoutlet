import React from 'react';
import List from '@material-ui/core/List';

import GroupListItem from './GroupListItem';
import OutletListItem from './OutletListItem';
import { apiRequest } from '../util';

class Group extends React.Component {
  state = {
    outlets: [],
  }

  componentDidMount() {
    const { outlets } = this.props;

    this.setState({ outlets });
  }

  handleAction = action => event => {
    const { id } = this.props;

    apiRequest('POST', '/outlet_group', { id, action })
      .then(result => result.outlets)
      .then(outlets => this.setState({ outlets }))
      .catch(err => console.error(err));
  }

  render() {
    const { name } = this.props;
    const { outlets } = this.state;

    return (
      <List component="nav">
        <GroupListItem
          name={name}
          onActionOn={this.handleAction('on')}
          onActionOff={this.handleAction('off')}
          onActionToggle={this.handleAction('toggle')}
        />
        {outlets.map(outlet =>
          <OutletListItem key={outlet.id} {...outlet} />
        )}
      </List>
    );
  }
}

export default Group;
