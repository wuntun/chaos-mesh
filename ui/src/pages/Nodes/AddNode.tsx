import { Box, Button, MenuItem } from '@material-ui/core'
import { Form, Formik, FormikHelpers } from 'formik'
import { SelectField, TextField } from 'components/FormField'

import { Node } from 'api/nodes'
import React from 'react'
import T from 'components/T'
import api from 'api'

type FormValues = Node

interface Props {
  onSubmitCallback?: (values: FormValues) => void
}

const AddNode: React.FC<Props> = ({ onSubmitCallback }) => {
  const submitNode = (values: FormValues, { resetForm }: FormikHelpers<FormValues>) => {
    api.nodes
      .add({
        ...values,
        config: window.btoa(values.config),
      })
      .then(() => {
        typeof onSubmitCallback === 'function' && onSubmitCallback(values)

        resetForm()
      })
  }

  return (
    <Formik initialValues={{ name: '', kind: 'k8s', config: '' }} onSubmit={submitNode}>
      <Form>
        <TextField id="name" name="name" label={T('nodes.add.name')} helperText={T('nodes.add.nameHelper')} />
        <SelectField id="kind" name="kind" label={T('nodes.add.kind')} helperText={T('nodes.add.kindHelper')}>
          <MenuItem value="k8s">Kubernetes</MenuItem>
          <MenuItem value="physic">Physic</MenuItem>
        </SelectField>
        <TextField
          id="config"
          name="config"
          label={T('nodes.add.config')}
          helperText={T('nodes.add.configHelper')}
          multiline
          rows={12}
        />
        <Box textAlign="right">
          <Button type="submit" variant="contained" color="primary">
            {T('common.submit')}
          </Button>
        </Box>
      </Form>
    </Formik>
  )
}

export default AddNode
