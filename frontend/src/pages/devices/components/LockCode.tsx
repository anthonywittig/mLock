import React from "react"
import { Badge, Button, Form } from "react-bootstrap"
import { addDays, format, parseISO, set } from "date-fns"
import { Loading } from "../../utils/Loading"
import { StandardFetch } from "../../utils/FetchHelper"

interface Props {
  deviceId: string
  managedLockCode: DeviceManagedLockCodeT | null
  managedLockCodesUpdated: () => void
}

export const LockCode = (props: Props) => {
  const [loading, setLoading] = React.useState<boolean>(false)
  let initialCode = ""
  if (props.managedLockCode) {
    initialCode =
      props.managedLockCode.code + " - " + props.managedLockCode.status
    if (props.managedLockCode.reservation.id) {
      initialCode +=
        " - Reservation " +
        props.managedLockCode.reservation.id.replace("@LiveRez.com", "")
    }
    initialCode += " - " + props.managedLockCode.note
  }
  const [code, setCode] = React.useState<string>(initialCode)
  const [syncWithReservation, setSyncWithReservation] = React.useState<boolean>(
    !!props.managedLockCode?.reservation.sync,
  )
  const [startAt, setStartAt] = React.useState<Date>(
    props.managedLockCode
      ? parseISO(props.managedLockCode.startAt)
      : new Date(),
  )
  const [endAt, setEndAt] = React.useState<Date>(
    props.managedLockCode
      ? parseISO(props.managedLockCode.endAt)
      : addDays(set(new Date(), { minutes: 0 }), 1),
  )

  const formSubmit = (evt: React.FormEvent<HTMLFormElement>) => {
    evt.preventDefault()
    setLoading(true)

    if (props.managedLockCode === null) {
      StandardFetch("devices/" + props.deviceId + "/lock-codes/", {
        method: "POST",
        body: JSON.stringify({
          code: code,
          startAt: startAt,
          endAt: endAt,
        }),
      })
        .then((response) => response.json())
        .then((response) => {
          setCode("")
          props.managedLockCodesUpdated()
          setLoading(false)
        })
        .catch((err) => {
          // TODO: indicate error.
          console.log(err)
        })
      return
    }

    StandardFetch(
      "devices/" + props.deviceId + "/lock-codes/" + props.managedLockCode?.id,
      {
        method: "PUT",
        body: JSON.stringify({
          reservation: {
            sync: syncWithReservation,
          },
          endAt: endAt,
          startAt: startAt,
        }),
      },
    )
      .then((response) => response.json())
      .then((response) => {
        props.managedLockCodesUpdated()
        setLoading(false)
      })
      .catch((err) => {
        // TODO: indicate error.
        console.log(err)
      })
  }

  const getDisabled = () => {
    if (
      props.managedLockCode?.status === "Adding" ||
      props.managedLockCode?.status === "Enabled"
    ) {
      return true
    }
    return syncWithReservation
  }

  const render = () => {
    // Might have just a "renderLockCode" method in the future too?
    return renderAddLockCode()
  }

  const renderAddLockCode = () => {
    if (loading) {
      return <Loading />
    }

    let statusBadge = <></>
    const status = props.managedLockCode?.status
    if (status) {
      const variant = (() => {
        switch (status) {
          case "Scheduled":
            return "warning"
          case "Adding":
            return "danger"
          case "Enabled":
            return "success"
          case "Removing":
            return "danger"
          case "Complete":
            return "secondary"
          default:
            return "danger"
        }
      })()

      statusBadge = (
        <Badge pill bg={variant}>
          {status}
        </Badge>
      )
    }

    let syncWithReservationCB = <></>
    if (props.managedLockCode?.reservation.id) {
      syncWithReservationCB = (
        <Form.Group className="mb-3" controlId="formBasicCheckbox">
          <Form.Check
            type="checkbox"
            label="Sync With Reservation"
            checked={syncWithReservation}
            onChange={(e) => setSyncWithReservation(e.target.checked)}
          />
        </Form.Group>
      )
    }

    // TODO: we should show the time zone that's being used.
    return (
      <Form onSubmit={(evt) => formSubmit(evt)} style={{ marginBottom: "2em" }}>
        <Form.Group>
          <Form.Label>Code {statusBadge}</Form.Label>
          <Form.Control
            type="text"
            defaultValue={code}
            onChange={(evt) => setCode(evt.target.value)}
            disabled={!!props.managedLockCode}
          />
        </Form.Group>

        {syncWithReservationCB}

        <Form.Group>
          <Form.Label>Enable At</Form.Label>
          <Form.Control
            type="datetime-local"
            defaultValue={format(startAt, "yyyy-MM-dd'T'HH:mm")}
            onChange={(evt) => setStartAt(parseISO(evt.target.value))}
            disabled={getDisabled()}
          />
        </Form.Group>

        <Form.Group>
          <Form.Label>Disable At</Form.Label>
          <Form.Control
            type="datetime-local"
            defaultValue={format(endAt, "yyyy-MM-dd'T'HH:mm")}
            onChange={(evt) => setEndAt(parseISO(evt.target.value))}
            disabled={syncWithReservation}
          />
        </Form.Group>

        <Button variant="secondary" type="submit">
          {props.managedLockCode ? "Update" : "Add"} Lock Code
        </Button>
      </Form>
    )
  }

  return render()
}
