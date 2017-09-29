import React from 'react';

import OutletGroup from './OutletGroup';

class OutletGroupContainer extends React.Component {
  render() {
    return (
      <div style={this.props.style}>
        {this.props.outletGroups.map((group, key) =>
          <OutletGroup
            key={key}
            identifier={group.identifier}
            outlets={group.outlets}
          />
        )}
      </div>
    );
  }
}

export default OutletGroupContainer;
