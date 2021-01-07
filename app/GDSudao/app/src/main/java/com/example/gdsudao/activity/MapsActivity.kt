package com.example.gdsudao.activity

import android.content.Intent
import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.widget.Toast
import com.example.gdsudao.R

import com.google.android.gms.maps.CameraUpdateFactory
import com.google.android.gms.maps.GoogleMap
import com.google.android.gms.maps.OnMapReadyCallback
import com.google.android.gms.maps.SupportMapFragment
import com.google.android.gms.maps.model.LatLng
import com.google.android.gms.maps.model.MarkerOptions
import kotlinx.android.synthetic.main.activity_maps.*

class MapsActivity : AppCompatActivity(),  OnMapReadyCallback {

    private lateinit var mMap: GoogleMap

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_maps)
        // Obtain the SupportMapFragment and get notified when the map is ready to be used.
        val mapFragment = supportFragmentManager
            .findFragmentById(R.id.map) as SupportMapFragment
        mapFragment.getMapAsync(this)

        val view = mapFragment.view
        view?.isClickable = false

    }

    /**
     * Manipulates the map once available.
     * This callback is triggered when the map is ready to be used.
     * This is where we can add markers or lines, add listeners or move the camera. In this case,
     * we just add a marker near Sydney, Australia.
     * If Google Play services is not installed on the device, the user will be prompted to install
     * it inside the SupportMapFragment. This method will only be triggered once the user has
     * installed Google Play services and returned to the app.
     */
    override fun onMapReady(googleMap: GoogleMap) {
        var contLocation = 0
        var latitude = ""
        var longitude = ""
        mMap = googleMap

        // Add a marker in Sydney and move the camera
        val epsul = LatLng(-31.3527, -54.0158)

        mMap.moveCamera(CameraUpdateFactory.newLatLng(epsul))

        mMap.setOnMapClickListener {
            //Toast.makeText(this, "tocou ${it.latitude} ${it.longitude} cont: ${contLocation}", Toast.LENGTH_SHORT).show()
            if (contLocation < 1){
                contLocation++
                val location = LatLng(it.latitude, it.longitude)
                latitude = it.latitude.toString()
                longitude = it.longitude.toString()
                val marker = mMap.addMarker(MarkerOptions().position(location).title("Localização atual"))
                marker.showInfoWindow()
            }

        }

        btnMapaReturnLocation.setOnClickListener {
            var intent = Intent(this, CadastroAreaActivity::class.java)
            if (latitude.isNotEmpty() && longitude.isNotEmpty()){
                intent.putExtra("lat", latitude)
                intent.putExtra("lon", longitude)
                startActivity(intent)
            }else{
                Toast.makeText(this, "Clique na tela pra selecionar um local", Toast.LENGTH_SHORT).show()
            }

        }
    }




}