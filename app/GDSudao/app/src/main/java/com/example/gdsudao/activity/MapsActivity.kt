package com.example.gdsudao.activity

import android.content.Intent
import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.util.Log
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

    var nomeArea = ""
    var dataCorte = ""
    var numeroCortes = ""
    var latitude = ""
    var longitude = ""

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_maps)
        // Obtain the SupportMapFragment and get notified when the map is ready to be used.
        val mapFragment = supportFragmentManager
            .findFragmentById(R.id.map) as SupportMapFragment
        mapFragment.getMapAsync(this)

        val view = mapFragment.view
        view?.isClickable = false


        val bundle = intent.extras

        if (bundle != null){
            nomeArea = bundle?.getString("nA").toString()
            dataCorte = bundle?.getString("dC").toString()
            numeroCortes = bundle?.getString("nC").toString()
            latitude = bundle?.getString("lat").toString()
            longitude = bundle?.getString("lon").toString()
            //Log.println(Log.DEBUG, "chegou?", "${nomeArea} + ${dataCorte} + ${numeroCortes} ")
        }


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
        mMap = googleMap

        // Add a marker in Sydney and move the camera
        val epsul = LatLng(-31.3527, -54.0158)

        mMap.moveCamera(CameraUpdateFactory.newLatLng(epsul))

        mMap.setOnMapClickListener {
            //Toast.makeText(this, "tocou ${it.latitude} ${it.longitude} cont: ${contLocation}", Toast.LENGTH_SHORT).show()
            if (contLocation < 1){
                contLocation++
                val location = LatLng(it.latitude, it.longitude)
                latitude = it.latitude.toString().substring(0, 11)
                longitude = it.longitude.toString().substring(0, 11)
                val marker = mMap.addMarker(MarkerOptions().position(location).title("Localização atual"))
                marker.showInfoWindow()
            }

        }

        btnMapaReturnLocation.setOnClickListener {

            if (latitude.isNotEmpty() && longitude.isNotEmpty()){
                //Log.println(Log.DEBUG, "localizacao", "${latitude} + ${longitude}")
                var intent = Intent(this, CadastroAreaActivity::class.java)
                intent.putExtra("lat", latitude)
                intent.putExtra("lon", longitude)
                intent.putExtra("nA", nomeArea)
                intent.putExtra("dC", dataCorte)
                intent.putExtra("nC", numeroCortes)
                startActivity(intent)
            }else{
                Toast.makeText(this, "Clique na tela pra selecionar um local", Toast.LENGTH_SHORT).show()
            }

        }
    }




}