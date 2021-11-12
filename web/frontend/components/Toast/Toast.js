import React from 'react';
import { toast } from 'react-toastify';
import Toast from 'react-bootstrap/Toast';

import 'react-toastify/dist/ReactToastify.min.css';

toast.configure({ hideProgressBar: true, closeButton: false });

export const showSuccessToast = ({ title = 'Success!', content = '' }) =>
  toast(({ closeToast }) => (
    <Toast onClose={closeToast}>
      <Toast.Header>
        <img
          src="/static/frontend/assets/images/trento-icon.png"
          className="rounded mr-2 toast-logo"
          alt=""
        />
        <strong className="mr-auto">{title}</strong>
      </Toast.Header>
      <Toast.Body className="trento-toast-body">{content}</Toast.Body>
    </Toast>
  ));

export const showErrorToast = ({ title = 'OH NOES', content = '' }) =>
  toast(({ closeToast }) => (
    <Toast onClose={closeToast}>
      <Toast.Header>
        <i className="eos-icons eos-18 color-red">error</i>
        <strong className="mr-auto">{title}</strong>
      </Toast.Header>
      <Toast.Body className="trento-toast-body">{content}</Toast.Body>
    </Toast>
  ));
