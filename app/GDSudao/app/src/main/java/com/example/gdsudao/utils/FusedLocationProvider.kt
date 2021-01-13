package com.example.gdsudao.utils

import android.content.Context
import android.location.Location
import com.google.android.gms.location.*

/**
 * A default implementation using the fused location provider to supply the application with locations
 *<p>
 * You should ensure you have the correct permissions in the manifest, that the permissions have
 * been granted, and location services are enabled.
 */
@SuppressWarnings("MissingPermission")
class FusedLocationProvider(ctx: Context) : LocationProvider {
    companion object {
        //How fast can your application process location updates requested by other applications
        private val DEFAULT_FASTEST_INTERVAL: Long = 1000 // 1 Second
        //Interval to ask for location updates until one has been determined
        private val DEFAULT_SHORT_INTERVAL: Long = 500 // 1/2 Second
        //Interval that is used AFTER a non-null location is received
        private val DEFAULT_LONG_INTERVAL: Long = 20 * 60 * 1000 //20 Minutes
        //The accuracy of the location updates
        private val LOCATION_PRIORITY = LocationRequest.PRIORITY_BALANCED_POWER_ACCURACY

        private fun isSameLatLong(l1: Location, l2: Location): Boolean = (l1.latitude == l2.latitude && l1.longitude == l2.longitude)
    }

    private val client = LocationServices.getFusedLocationProviderClient(ctx)
    private val locationCallbackHelper = LocationCallbackHelper()
    private var lastKnownLocation: Location? = null
    private var listener: LocationProvider.Listener? = null
    private var locationRequest = LocationRequest.create()

    override fun getCurrentLocation(): Location? = lastKnownLocation

    override fun setListener(listener: LocationProvider.Listener?) {
        this.listener = listener
    }

    override fun startUpdates() {
        stopUpdates()
        client.lastLocation
            .addOnSuccessListener { location -> handleLocationUpdate(location) }
            .addOnCompleteListener {
                val interval = if (lastKnownLocation == null) DEFAULT_SHORT_INTERVAL else DEFAULT_LONG_INTERVAL
                requestLocationUpdates(interval, LOCATION_PRIORITY)
            }
    }

    override fun stopUpdates() {
        client.removeLocationUpdates(locationCallbackHelper)
        reset()
    }

    //Permissions and location services should be enabled by this point, make sure you do it!!
    private fun requestLocationUpdates(interval: Long, priority: Int) {
        locationRequest.setFastestInterval(DEFAULT_FASTEST_INTERVAL)
            .setInterval(interval)
            .setPriority(priority)
        client.requestLocationUpdates(locationRequest, locationCallbackHelper, null)
    }

    private fun handleLocationUpdate(location: Location?) {
        val lastLocation = lastKnownLocation
        lastKnownLocation = location

        //Update to a longer interval, also will trigger a location update using the last one
        if (location != null) {
            if (lastLocation == null || !isSameLatLong(lastLocation, location)) listener?.onNewLocationUpdate(location)
            if (locationRequest.interval == DEFAULT_SHORT_INTERVAL) {
                requestLocationUpdates(DEFAULT_LONG_INTERVAL, LOCATION_PRIORITY)
                return
            }
        }
        listener?.onLocationUpdate(location)
    }

    private fun reset() {
        lastKnownLocation = null
    }

    private inner class LocationCallbackHelper : LocationCallback() {
        override fun onLocationResult(lr: LocationResult) {
            handleLocationUpdate(lr.lastLocation)
        }

        override fun onLocationAvailability(la: LocationAvailability) {
            if (lastKnownLocation == null && la.isLocationAvailable) {
                client.lastLocation.addOnSuccessListener { location ->
                    handleLocationUpdate(location)
                }
            }
        }
    }
}