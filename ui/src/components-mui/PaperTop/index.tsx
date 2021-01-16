import { Box, Typography } from '@material-ui/core'

import React from 'react'
import { makeStyles } from '@material-ui/core/styles'

const useStyles = makeStyles({
  root: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    width: '100%',
    height: 56,
  },
  title: {
    fontWeight: 700,
  },
})

interface PaperTopProps {
  title?: string | JSX.Element
  subtitle?: string | JSX.Element
}

const PaperTop: React.FC<PaperTopProps> = ({ title, subtitle, children }) => {
  const classes = useStyles()

  return (
    <Box className={classes.root} px={3}>
      <Box>
        <Typography className={classes.title}>{title}</Typography>
        <Typography variant="body2" color="textSecondary">
          {subtitle}
        </Typography>
      </Box>
      {children}
    </Box>
  )
}

export default PaperTop
