import React from 'react';
import PropTypes from 'prop-types';
import Dialog from '@material-ui/core/Dialog';

import DialogAppBar from './DialogAppBar';

export default function ConfigurationDialog(props) {
  const {
    children,
    open,
    onClose,
    onDone,
    doneButtonDisabled,
    doneButtonText,
    title
  } = props;

  return (
    <Dialog fullScreen open={open} onClose={onClose}>
      <DialogAppBar
        title={title}
        onClose={onClose}
        onDone={onDone}
        doneButtonDisabled={doneButtonDisabled}
        doneButtonText={doneButtonText}
      />
      {children}
    </Dialog>
  );
}

ConfigurationDialog.propTypes = {
  open: PropTypes.bool,
  onClose: PropTypes.func.isRequired,
  onDone: PropTypes.func,
  title: PropTypes.string.isRequired,
  doneButtonDisabled: PropTypes.bool,
  doneButtonText: PropTypes.string,
};
